use super::session::Session;
use crate::common::Group;
use crate::error::{ErrorKind, Group as GroupError};

impl Session {
    pub fn create_group(&mut self, group_name: &str) -> Result<Group, ErrorKind> {
        if self.exist_group(group_name) {
            return Err(ErrorKind::Group(GroupError::Exist(group_name.to_string())));
        }

        let group_id = self.groups.len() as i32 + 1;
        let group_uuid = uuid::Uuid::new_v4();

        let group = Group {
            group_uuid,
            group_id,
            group_name: group_name.to_string(),
        };

        self.groups.push(group.clone());
        Ok(group)
    }

    pub fn get_group(&self, group_name: &str) -> Result<Group, ErrorKind> {
        for group in &self.groups {
            if group.group_name.eq(group_name) {
                return Ok(group.clone());
            }
        }

        Err(ErrorKind::Group(GroupError::Get(group_name.to_string())))
    }

    pub fn get_group_by_id(&self, group_id: &str) -> Result<Group, String> {
        for group in &self.groups {
            if group.group_id.to_string().eq(group_id) {
                return Ok(group.clone());
            }
        }

        Err("".to_string())
    }

    /// checking if group exist with group_name
    pub fn exist_group(&self, name: &str) -> bool {
        for group in &self.groups {
            if group.group_name.eq(name) {
                return true;
            }
        }
        false
    }

    /// checking if group_id exist in current session
    pub fn exist_group_id(&self, id: &str) -> bool {
        for group in &self.groups {
            if group.group_id.to_string().eq(id) {
                return true;
            }
        }
        false
    }
}
