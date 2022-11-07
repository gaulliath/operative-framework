use crate::common as opf_common;
use serde::{Deserialize, Serialize};
use std::time::SystemTime;
use uuid;

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct Session {
    pub session_id: uuid::Uuid,
    pub session_name: String,
    pub targets: Vec<opf_common::Target>,
    pub groups: Vec<opf_common::Group>,
    pub links: Vec<opf_common::Link>,
    pub actions: Vec<opf_common::Action>,
    pub created_at: SystemTime,
}

#[derive(Debug, Default)]
pub struct SessionConfig {
    pub verbose: bool,
    pub temporary: bool,
}

pub fn new() -> (Session, SessionConfig) {
    (
        Session {
            session_id: uuid::Uuid::new_v4(),
            session_name: String::new(),
            targets: vec![],
            groups: vec![],
            links: vec![],
            actions: vec![],
            //      vm: rlua::Lua::new(),
            created_at: SystemTime::now(),
        },
        SessionConfig::default(),
    )
}
