use thiserror::Error;

#[derive(Error, Debug)]
pub enum Error {
    #[error("module not found")]
    ModuleNameNotFound,
    #[error("group not found")]
    GroupNotFound,
    #[error("cant load content of module")]
    CantLoadContent,
    #[error("parameter {0} not specified, please set")]
    ParamNotAvailable(String),
    #[error("target not available in current session")]
    TargetNotAvailable,
}
