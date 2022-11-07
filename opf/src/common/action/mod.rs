use serde::{Deserialize, Serialize};
use strum_macros::{Display, EnumString};
use thiserror::Error;

#[derive(Debug, PartialEq, EnumString, Display, Clone, Serialize, Deserialize)]
#[strum(serialize_all = "lowercase")]
pub enum ActionType {
    MergeTarget,
}

#[derive(Debug, PartialEq, EnumString, Display, Clone, Serialize, Deserialize)]
#[strum(serialize_all = "lowercase")]
pub enum ActionState {
    Pending,
    Complete,
    Canceled,
}

#[derive(Error, Debug)]
pub enum Error {
    #[error("action {0} is not available")]
    ActionNotFound(String),
}
