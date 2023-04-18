use std::collections::HashMap;
use std::str::FromStr;
use std::time::SystemTime;
use serde::{Deserialize, Serialize};
use strum_macros::{Display, EnumString, EnumIter};

use error::ErrorKind;

pub mod error;
pub mod event;
mod impl_model;
pub mod link;
pub mod metadata;
pub mod module;
pub mod target;
pub mod workspace;
pub mod command;

#[derive(Debug)]
pub struct Command {
    pub action: CommandAction,
    pub object: CommandObject,
    pub params: HashMap<String, String>,
}

#[derive(Debug, PartialEq, EnumString)]
#[strum(serialize_all = "lowercase")]
pub enum CommandAction {
    Add,
    Del,
    List,
    Set,
    Run,
    Accept,
    Save,
    Load,
    Export,
    Connect,
    Switch,
    Stop,
    Help,
}

#[derive(Debug, PartialEq, EnumString)]
#[strum(serialize_all = "lowercase")]
pub enum CommandObject {
    None,
    Save,
    Target,
    Link,
    Session,
    Action,
    Group,
    Api,
    Workspace,
    Export(String),
    Module(String),
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Target {
    pub target_id: i32,
    pub target_type: TargetType,
    pub target_name: String,
    pub target_groups: Vec<i32>,
    pub target_custom_id: Option<String>,
    pub target_parent: Option<i32>,
    pub meta: HashMap<String, String>,
}

#[derive(Debug, PartialEq, EnumString, EnumIter, Display, Clone, Serialize, Deserialize)]
#[strum(serialize_all = "lowercase")]
pub enum TargetType {
    Person,
    Company,
    Alias,
    Email,
    Username,
    Port,
    Document,
    PhoneNumber,
    Image,
    Account,
    IpAddress,
    Domain,
}

#[derive(Debug, PartialEq, EnumString, Display, Clone, Serialize, Deserialize)]
#[strum(serialize_all = "lowercase")]
pub enum LinkType {
    Both,
    In,
    Out,
}

#[derive(Debug, PartialEq, EnumString, Display, Clone, Serialize, Deserialize)]
#[strum(serialize_all = "lowercase")]
pub enum LinkFrom {
    CLI,
    Module(String),
    Target(i32),
    Other,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Group {
    pub group_id: i32,
    pub group_name: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Link {
    pub link_id: i32,
    pub link_type: LinkType,
    pub link_meta: HashMap<String, String>,
    pub link_from: i32,
    pub link_to: i32,
    pub link_created_by: LinkFrom,
    pub link_created_at: SystemTime,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Workspace {
    pub workspace_id: i32,
    pub workspace_name: String,
    pub workspace_created_at: SystemTime,
}

pub fn validate_type(t: &str) -> Result<TargetType, ErrorKind> {
    match TargetType::from_str(t.to_lowercase().as_str()) {
        Ok(t) => Ok(t),
        Err(_) => Err(error::ErrorKind::Target(error::Target::ParamTypeNotFound)),
    }
}

pub fn validate_link_type(t: &str) -> Result<LinkType, ErrorKind> {
    match LinkType::from_str(t.to_lowercase().as_str()) {
        Ok(t) => Ok(t),
        Err(_) => Err(ErrorKind::Link(error::Link::ParamTypeNotFound)),
    }
}

pub fn validate_link_created_by(t: &str) -> Result<LinkFrom, ErrorKind> {
    match LinkFrom::from_str(t.to_lowercase().as_str()) {
        Ok(t) => Ok(t),
        Err(_) => Err(ErrorKind::Link(error::Link::ParamTypeNotFound)),
    }
}

impl Target {
    pub fn to_hashmap(&mut self) -> HashMap<String, String> {
        let mut map = HashMap::new();
        map.insert(target::ID.to_string(), self.target_id.to_string());
        map.insert(target::TYPE.to_string(), self.target_type.to_string());
        map.insert(target::NAME.to_string(), self.target_name.to_string());
        map.insert(
            target::CUSTOM_ID.to_string(),
            self.target_custom_id.clone().unwrap_or(String::new()),
        );
        map.insert(
            target::PARENT.to_string(),
            self.target_parent.unwrap_or(0).to_string(),
        );
        self.meta.remove(target::ID);
        self.meta.remove(target::PARENT);
        self.meta.remove(target::TYPE);
        self.meta.remove(target::NAME);

        map.extend(self.meta.clone());
        map
    }
}

impl Link {
    pub fn to_hashmap(&mut self) -> HashMap<String, String> {
        let mut map = HashMap::new();
        map.insert(link::ID.to_string(), self.link_id.to_string());
        map.insert(link::TYPE.to_string(), self.link_type.to_string());
        map.insert(link::TO.to_string(), self.link_to.to_string());
        map.insert(link::FROM.to_string(), self.link_from.to_string());
        map.insert(
            link::CREATED_BY.to_string(),
            self.link_created_by.to_string(),
        );
        map.extend(self.link_meta.clone());
        map
    }
}
