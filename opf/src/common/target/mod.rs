use crate::error::{ErrorKind, Target as TargetError};
use serde::{Deserialize, Serialize};
use std::str::FromStr;
use strum_macros::{Display, EnumString};
use thiserror::Error;

#[derive(Debug, PartialEq, EnumString, Display, Clone, Serialize, Deserialize)]
#[strum(serialize_all = "lowercase")]
pub enum TargetType {
    Person,
    Company,
    Alias,
    Email,
    Document,
    PhoneNumber,
    Image,
    Account,
    IpAddress,
    Domain,
}

#[derive(Error, Debug)]
pub enum Error {
    #[error("Parameter 'name' mandatory for this action")]
    ParamNameNotFound,
    #[error("Parameter 'type' mandatory for this action")]
    ParamTypeNotFound,
    #[error("Target already found in this session")]
    TargetExist,
    #[error("Target type '{0}' not available")]
    TypeNotAvailable(String),
    #[error("Parent not found")]
    ParentUuidNotFound,
    #[error("Parent uuid ins't valid")]
    ParentUuidNotValid,
}

pub fn validate_type(t: &str) -> Result<TargetType, ErrorKind> {
    match TargetType::from_str(t.clone().to_lowercase().as_str()) {
        Ok(t) => Ok(t),
        Err(_) => Err(ErrorKind::Target(TargetError::ParamTypeNotFound)),
    }
}
