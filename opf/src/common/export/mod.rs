use strum_macros::{Display,EnumString};
use thiserror::Error;

#[derive(Error, Debug)]
pub enum Error {
    #[error("module not found")]
    ModuleNameNotFound,
    #[error("cant load content of module")]
    CantLoadContent,
    #[error("export type not found")]
    ExportType,
}

#[derive(Debug, Display, EnumString)]
#[strum(serialize_all = "lowercase")]
pub enum ExportType {
    Dot
}
