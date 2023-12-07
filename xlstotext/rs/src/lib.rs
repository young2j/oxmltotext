pub mod app;
pub mod errors;
pub mod totext;

use std::path::PathBuf;

pub use app::App;
pub use errors::Errors;
pub use totext::XlsToText;

use clap::{CommandFactory, FromArgMatches};

pub fn run(rs: &'static str, fs: &'static str, ss: &'static str) -> Result<(), Errors> {
    let app = App::from_arg_matches(
        &App::command()
            .long_about(&format!(
                "\
Converts the spreadsheet at PATH (or <sheet> if \
requested) to text sent to stdout.

Should be able to convert from (and automatically guess between) \
XLS, XLSX, XLSB and ODS."
            ))
            .mut_arg("row_separator", |arg| arg.default_value(rs))
            .mut_arg("col_separator", |arg| arg.default_value(fs))
            .mut_arg("sheet_separator", |arg| arg.default_value(ss))
            .get_matches(),
    )
    .unwrap();

    XlsToText::new(app).parse_worksheets()?;

    Ok(())
}

pub fn bench_run(path: &str) -> Result<(), Errors> {
    let app = App {
        path: PathBuf::from(path),
        sheet: "all".to_string(),
        row_separator: "\n".to_string(),
        col_separator: "\t".to_string(),
        sheet_separator: "-".to_string(),
        formula: app::FormulaMode::CachedValue,
    };

    XlsToText::new(app).parse_worksheets()?;

    Ok(())
}
