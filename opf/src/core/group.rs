use super::session::Session;
use crate::common::Group;

impl Session {
    pub fn create_group<'a>(&mut self, group_name: &'a String) -> Result<Group, String> {
        if self.exist_group(group_name) {
            return Err("".to_string());
        }

        let group_id = self.groups.len() as i32 + 1;
        let group_uuid = uuid::Uuid::new_v4();

        let group = Group {
            group_uuid,
            group_id,
            group_name: group_name.clone()
        };

        self.groups.push(group.clone());
        Ok(group)
    }

    pub fn get_group<'a>(&self, group_name: &'a str) -> Result<Group, String> {
        for group in &self.groups {
            if group.group_name.eq(group_name) {
                return Ok(group.clone())
            }
        }

        Err("".to_string())
    }

    pub fn get_group_by_id<'a>(&self, group_id: &'a str) -> Result<Group, String> {
        for group in &self.groups {
            if group.group_id.to_string().eq(group_id) {
                return Ok(group.clone())
            }
        }

        Err("".to_string())
    }

    /// checking if group exist with group_name
    pub fn exist_group<'a>(&self, name: &'a String) -> bool {
        for group in &self.groups {
            if group.group_name.eq(name) {
                return true;
            }
        }
        false
    }

    /// checking if group_id exist in current session
    pub fn exist_group_id<'a>(&self, id: &'a String) -> bool {
        for group in &self.groups {
            if group.group_id.to_string().eq(id) {
                return true;
            }
        }
        false
    }
}
