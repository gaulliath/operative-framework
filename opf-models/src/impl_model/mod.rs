use std::collections::HashMap;
use std::time::SystemTime;

use crate::error::{self, ErrorKind};
use crate::{link, target, validate_link_created_by, validate_link_type, validate_type};

impl crate::Command {
    pub fn get<'a>(&self, name: &str) -> Option<String> {
        if let Some(value) = self.params.get(&name.to_string()) {
            return Some(value.clone());
        }
        None
    }
}

impl TryFrom<HashMap<String, String>> for crate::Target {
    type Error = ErrorKind;

    fn try_from(mut params: HashMap<String, String>) -> Result<Self, Self::Error> {
        let target_name = match params.get(target::NAME) {
            Some(name) => name.clone(),
            None => {
                return Err(ErrorKind::Target(error::Target::ParamNameNotFound));
            }
        };

        let target_id = match params.get(target::ID) {
            Some(id) => id.parse::<i32>().unwrap_or(0),
            None => 0,
        };

        let target_type = match params.get(target::TYPE) {
            Some(t) => match validate_type(t) {
                Ok(t) => t,
                Err(e) => return Err(ErrorKind::Target(error::Target::Parsing(e.to_string()))),
            },
            None => {
                return Err(ErrorKind::Target(error::Target::ParamTypeNotFound));
            }
        };

        let target_parent = match params.get(target::PARENT) {
            Some(parent_id) => Some(
                parent_id
                    .parse::<i32>()
                    .map_err(|e| ErrorKind::GenericError(e.to_string()))?,
            ),
            None => None,
        };

        let target_custom_id = match params.get(target::CUSTOM_ID) {
            Some(custom_id) => Some(custom_id.clone()),
            None => None,
        };

        params.remove(target::NAME);
        params.remove(target::TYPE);
        params.remove(target::CUSTOM_ID);

        Ok(Self {
            target_id,
            target_name,
            target_type,
            target_groups: vec![],
            target_custom_id,
            target_parent,
            meta: params.clone(),
        })
    }
}

impl TryFrom<HashMap<String, String>> for crate::Link {
    type Error = ErrorKind;

    fn try_from(mut params: HashMap<String, String>) -> Result<Self, Self::Error> {
        let link_id = match params.get(link::ID) {
            Some(id) => id.parse::<i32>().unwrap_or(0),
            None => 0,
        };

        let link_type = params
            .get(link::TYPE)
            .ok_or(ErrorKind::Link(error::Link::ParamTypeNotFound))
            .map(|t| validate_link_type(t))??;

        let link_created_by = params
            .get(link::CREATED_BY)
            .ok_or(ErrorKind::Link(error::Link::ParamTypeNotFound))
            .map(|t| validate_link_created_by(t))??;

        let link_to = params
            .get(link::TO)
            .ok_or(ErrorKind::Link(error::Link::ParamToNotFound))
            .map(|id| {
                id.parse::<i32>()
                    .map_err(|_| ErrorKind::InvalidFormatArgument)
            })??;

        let link_from = params
            .get(link::FROM)
            .ok_or(ErrorKind::Link(error::Link::ParamFromNotFound))
            .map(|id| {
                id.parse::<i32>()
                    .map_err(|_| ErrorKind::InvalidFormatArgument)
            })??;

        params.remove(link::ID);
        params.remove(link::TO);
        params.remove(link::FROM);
        params.remove(link::TYPE);
        params.remove(link::CREATED_BY);

        Ok(Self {
            link_id,
            link_to,
            link_created_by,
            link_from,
            link_type,
            link_meta: params.clone(),
            link_created_at: SystemTime::now(),
        })
    }
}
