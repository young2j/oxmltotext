use xlstotext::{run, Errors};

fn main() -> Result<(), Errors> {
    run("\n", "\t", "-")
}
