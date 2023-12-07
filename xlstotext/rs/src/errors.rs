use std::error::Error;
use std::fmt::{self, Debug, Display, Formatter};

pub enum Errors {
    InvalidSeparator,
    MissingSeparator,
    Empty,
    NotFound(String),
    Csv(csv::Error),
    Spreadsheet(calamine::Error),
    CellError(calamine::CellErrorType),
}

impl Error for Errors {}
// delegate to Display so the error message is not crap
impl Debug for Errors {
    fn fmt(&self, fmt: &mut Formatter) -> fmt::Result {
        (&self as &dyn Display).fmt(fmt)
    }
}
impl Display for Errors {
    fn fmt(&self, fmt: &mut Formatter) -> fmt::Result {
        use Errors::*;
        match self {
            Empty => write!(fmt, "Empty spreadsheet"),
            NotFound(s) => write!(fmt, "Could not find sheet {s:?} in spreadsheet"),
            InvalidSeparator => write!(
                fmt,
                "A provided separator is invalid, separators need to be a single ascii chacter"
            ),
            MissingSeparator => write!(fmt, "No separator found"),
            Csv(e) => write!(fmt, "{e}"),
            Spreadsheet(e) => write!(fmt, "{e}"),
            CellError(e) => write!(fmt, "Error found in cell ({e:?})"),
        }
    }
}

impl From<calamine::Error> for Errors {
    fn from(err: calamine::Error) -> Self {
        Self::Spreadsheet(err)
    }
}

impl From<csv::Error> for Errors {
    fn from(err: csv::Error) -> Self {
        Self::Csv(err)
    }
}
