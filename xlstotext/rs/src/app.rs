use clap::{Parser, ValueEnum};
use std::path::PathBuf;

#[derive(Copy, Clone, Debug, PartialEq, Eq, ValueEnum)]
pub enum FormulaMode {
    /// never show formula, always display cached value, even if empty
    CachedValue,
    /// show formula if cached value is empty or absent
    IfEmpty,
    /// always show formula for formula cells (ignore cached values)
    Always,
    // TODO: evaluate formulas
    // Evaluate,
}

#[derive(Parser, Debug)]
#[command(author, version)]
#[command(about = "Converts spreadsheets to text")]
pub struct App {
    /// Spreadsheet file path
    pub path: PathBuf,

    /// Names or indexes (1 is first, separated by comma) of the sheet to convert
    #[arg(short, long, default_value = "all")]
    pub sheet: String,

    /// Row separator (a single character)
    #[arg(short, long, default_value = "\n")]
    pub row_separator: String,

    /// Column separator (a single character)
    #[arg(short, long, default_value = "\t")]
    pub col_separator: String,

    /// Sheet separator (a single character) x100
    #[arg(long, default_value = "-")]
    pub sheet_separator: String,

    /// Whether and when to show formulas
    #[arg(long, value_enum, default_value_t = FormulaMode::CachedValue)]
    pub formula: FormulaMode,
}
