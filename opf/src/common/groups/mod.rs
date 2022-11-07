use thiserror::Error;

#[derive(Error, Debug)]
pub enum Error {
    #[error("{0}")]
    GenericError(String),
    #[error("link already found in this session")]
    LinkExist,
}
