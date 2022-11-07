pub mod target;
pub mod search;
pub mod link;
pub mod module;
pub mod export;
pub mod action;
pub mod groups;

use std::collections::HashMap;
use std::time::SystemTime;
use serde::{Serialize, Deserialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Target {
    pub target_uuid: uuid::Uuid,
    pub target_id: i32,
    pub target_type: target::TargetType,
    pub target_name: String,
    pub target_groups: Vec<uuid::Uuid>,
    pub target_custom_id: Option<String>,
    pub target_parent: Option<uuid::Uuid>,
    pub meta: HashMap<String, String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Group {
    pub group_uuid: uuid::Uuid,
    pub group_id: i32,
    pub group_name: String,
}


#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Link {
    pub link_id: uuid::Uuid,
    pub link_label: String,
    pub link_color: Option<String>,
    pub link_source: uuid::Uuid,
    pub link_target: uuid::Uuid,
    pub link_created_by: String,
    pub link_created_at: SystemTime,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Action {
    pub action_id: i32,
    pub action_type: action::ActionType,
    pub action_ctx: Vec<u8>,
    pub action_state: action::ActionState,
    pub action_created_by: String,
    pub action_created_at: SystemTime,
}


impl TryFrom<HashMap<String, String>> for Target {
    type Error = String;

    fn try_from(mut params: HashMap<String, String>) -> Result<Self, Self::Error> {
        let target_name = match params.get("name") {
            Some(name) => name.clone(),
            None => {
                let e = target::Error::ParamNameNotFound.to_string();
                return Err(e);
            }
        };

        let target_type = match params.get("type") {
            Some(t) => {
                match target::validate_type(t.to_string()) {
                    Ok(t) => t,
                    Err(e) => {
                        return Err(e.to_string())
                    },
                }
            },
            None => {
                let e = target::Error::ParamTypeNotFound.to_string();
                return Err(e);
            }
        };

        let target_custom_id = match params.get("custom_id") {
            Some(custom_id) => Some(custom_id.clone()),
            None => None
        };

        params.remove("name");
        params.remove("type");
        params.remove("custom_id");



        Ok(Self {
            target_uuid: uuid::Uuid::new_v4(),
            target_id: 0,
            target_name,
            target_type,
            target_groups: vec![],
            target_custom_id,
            target_parent: None,
            meta: params.clone()
        })
    }
}


#[cfg(test)]
mod tests {
    #[test]
    fn it_works() {
        let result = 2 + 2;
        assert_eq!(result, 4);
    }
}