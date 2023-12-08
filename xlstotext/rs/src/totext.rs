use std::{
    fmt::Write,
    fs::File,
    io::{self, BufReader},
};

use calamine::{open_workbook_auto, DataType, Reader, Sheets};
use guard::guard;

use crate::{
    app::{App, FormulaMode},
    errors::Errors,
};

pub struct XlsToText {
    app: App,
    wb: Sheets<BufReader<File>>,
}

fn separator_to_byte(s: &str) -> Result<u8, Errors> {
    let c = s.chars().next().ok_or(Errors::InvalidSeparator)?;
    (c as u32).try_into().map_err(|_| Errors::InvalidSeparator)
}

impl XlsToText {
    pub fn new(app: App) -> XlsToText {
        let wb = open_workbook_auto(&app.path).unwrap();

        XlsToText { app, wb }
    }

    pub fn parse_worksheets(&mut self) -> Result<(), Errors> {
        // if sheet is a number get corresponding sheet in list, otherwise
        // assume all sheet names
        let sheet_names = self.wb.sheet_names();
        let mut parse_names = Vec::new();
        self.app.sheet.split(",").for_each(|s| {
            let _ = s.parse::<usize>().and_then(|n| {
                if let Some(name) = sheet_names.get(n.saturating_sub(1)) {
                    parse_names.push(name.to_string());
                };

                Ok(())
            });
        });

        if parse_names.len() > 0 {
            for name in parse_names {
                self.parse_worksheet(name)?;
            }
        } else {
            for name in sheet_names {
                self.parse_worksheet(name)?;
            }
        }

        Ok(())
    }

    fn parse_worksheet(&mut self, sheet_name: String) -> Result<(), Errors> {
        guard!(let Some(Ok(range)) = self.wb.worksheet_range(&sheet_name) else {
            return Err(Errors::NotFound(sheet_name));
        });
        guard!(let Some((offset_j, offset_i)) = range.start() else {
            return Ok(());
        });

        let formula_range = self
            .wb
            .worksheet_formula(&sheet_name)
            .expect("we know the sheet exists");
        let formatter: Box<dyn Fn(u32, u32, DataType) -> DataType> = match formula_range.as_ref() {
            Ok(f) => match self.app.formula {
                FormulaMode::CachedValue => Box::new(|_, _, cell| cell),
                FormulaMode::IfEmpty => Box::new(|i, j, cell| {
                    let formula = f.get_value((j, i)).filter(|s| !s.is_empty());
                    match cell {
                        DataType::Empty => {
                            formula.map_or(DataType::Empty, |v| DataType::String(v.to_string()))
                        }

                        DataType::String(s) if s.is_empty() => {
                            formula.map_or(DataType::Empty, |v| DataType::String(v.to_string()))
                        }

                        rest => rest,
                    }
                }),
                FormulaMode::Always => Box::new(|j, i, cell| {
                    f.get_value((i, j))
                        .filter(|s| !s.is_empty())
                        .map_or(cell, |s| DataType::String(s.to_string()))
                }),
            },
            Err(e) => {
                if self.app.formula != FormulaMode::CachedValue {
                    eprintln!("Formula parsing error: {e:?}");
                }
                Box::new(|_, _, cell| cell)
            }
        };

        let stdout = io::stdout().lock();
        let mut out = csv::WriterBuilder::new()
            .terminator(csv::Terminator::Any(separator_to_byte(
                &self.app.row_separator,
            )?))
            .delimiter(separator_to_byte(&self.app.col_separator)?)
            .from_writer(stdout);

        let mut contents = vec![String::new(); range.width()];
        for (j, row) in range.rows().enumerate() {
            for (i, (c, cell)) in row.iter().zip(contents.iter_mut()).enumerate() {
                cell.clear();
                match formatter(i as u32 + offset_i, j as u32 + offset_j, c.clone()) {
                    // DataType::Error(e) => return Err(Errors::CellError(e)),
                    // dont't interrupt other cells parsing
                    DataType::Error(e) => eprintln!("Error: {e}"),
                    // don't bother updating cell for empty
                    DataType::Empty => (),
                    // don't go through fmt for strings
                    DataType::String(s) => cell.push_str(&s),
                    rest => write!(cell, "{rest}")
                        .expect("formatting basic types to a string should never fail"),
                };
            }

            out.write_record(&contents)?;
        }

        let mut sheet_separator = vec![String::new(); range.width() - 1];
        sheet_separator.insert(0, self.app.sheet_separator.repeat(100));
        out.write_record(&sheet_separator)?;

        out.flush().unwrap();

        Ok(())
    }
}
