use thiserror::Error;

#[derive(Error, Debug)]
pub enum Error<'a> {
    #[error("parameter '{0}' mandatory for this action")]
    ParamNotFound(&'a str),
    #[error("format invalid for {0} parameter")]
    ParamFormatInvalid(&'a str),
    #[error("Parameter 'label' mandatory for this action")]
    ParamLabelNotFound,
    #[error("Parameter 'type' mandatory for this action")]
    ParamTypeNotFound,
    #[error("link already found in this session")]
    LinkExist,
}
