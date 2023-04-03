use super::target::TargetType;
use std::{collections::HashMap, str::FromStr};

#[derive(Debug, Default)]
pub struct Target {
    pub target_id: Option<String>,
    pub target_uuid: Option<String>,
    pub target_type: Option<TargetType>,
    pub target_name: Option<String>,
    pub target_parent: Option<String>,
    pub target_group_uuid: Option<String>,
    pub target_custom_id: Option<String>,
    pub meta: Option<HashMap<String, String>>,
}

impl From<HashMap<String, String>> for Target {
    fn from(mut e: HashMap<String, String>) -> Self {
        let target_id = match e.get("id") {
            Some(id) => Some(id.clone()),
            None => None,
        };

        let target_uuid = match e.get("uuid") {
            Some(uuid) => Some(uuid.clone()),
            None => None,
        };

        let target_type = match e.get("type") {
            Some(t) => match TargetType::from_str(t.clone().as_str()) {
                Ok(t) => Some(t),
                Err(_) => None,
            },
            None => None,
        };

        let target_parent = match e.get("parent") {
            Some(parent) => Some(parent.clone()),
            None => None,
        };

        let target_group_uuid = match e.get("group_uuid") {
            Some(group_uuid) => Some(group_uuid.clone()),
            None => None,
        };

        let target_name = match e.get("name") {
            Some(name) => Some(name.clone()),
            None => None,
        };

        let target_custom_id = match e.get("custom_id") {
            Some(cid) => Some(cid.clone()),
            None => None,
        };

        e.remove("id");
        e.remove("parent");
        e.remove("uuid");
        e.remove("type");
        e.remove("name");
        e.remove("group_uuid");

        Self {
            target_uuid,
            target_id,
            target_parent,
            target_type,
            target_name,
            target_group_uuid,
            target_custom_id,
            meta: Some(e),
        }
    }
}

#[derive(Debug)]
pub struct Link {
    pub link_id: Option<String>,
    pub link_label: Option<String>,
    pub link_color: Option<String>,
    pub link_source: Option<String>,
    pub link_target: Option<String>,
    pub link_created_by: Option<String>,
}

impl From<HashMap<String, String>> for Link {
    fn from(e: HashMap<String, String>) -> Self {
        let link_id = match e.get("id") {
            Some(id) => Some(id.clone()),
            None => None,
        };

        let link_source = match e.get("source") {
            Some(source) => Some(source.clone()),
            None => None,
        };

        let link_target = match e.get("target") {
            Some(target) => Some(target.clone()),
            None => None,
        };

        let link_label = match e.get("label") {
            Some(label) => Some(label.clone()),
            None => None,
        };

        let link_color = match e.get("color") {
            Some(color) => Some(color.clone()),
            None => None,
        };

        Self {
            link_id,
            link_label,
            link_color,
            link_source,
            link_target,
            link_created_by: None,
        }
    }
}
