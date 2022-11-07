use thiserror::Error;

#[derive(Error, Debug)]
pub enum Manager {
    #[error("can't process file")]
    CantProcessFile,
    #[error("directory {0} cant be opened")]
    CantOpenDirectory(String),
    #[error("can't read a content of module {0}")]
    CantReadContent(String),
    #[error("parsing metadata of {0} is impossible, please check format")]
    CantParseMetadata(String)
}
