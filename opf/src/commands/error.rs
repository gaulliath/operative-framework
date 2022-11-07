use thiserror::Error;

#[derive(Error, Debug)]
pub enum Error {
    #[error("argument is not found for this command")]
    InvalidCommandArgument,
    #[error("invalid format for argument")]
    InvalidFormatArgument,
    #[error("action not available")]
    ActionNotAvailable,
    #[error("context object for this action not available")]
    ObjectNotAvailable,
    #[error("{0}")]
    GenericError(String),
}
