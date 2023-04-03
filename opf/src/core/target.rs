use std::str::FromStr;

use super::session::Session;
use crate::common::search::Target as TargetSearch;
use crate::common::{self, Group, Target};
use crate::error::{ErrorKind, Target as TargetError};

impl Session {
    /// create a new target
    pub fn create_target(&mut self, mut target: common::Target) -> Result<Target, ErrorKind> {
        let mut target_search = TargetSearch::default();
        target_search.target_name = Some(target.target_name.clone());
        target_search.target_type = Some(target.target_type.clone());

        if self.exist_target(&target_search) {
            return Err(ErrorKind::Target(TargetError::CantBeCreated));
        }

        let target_id = self.targets.len() as i32 + 1;
        target.target_id = target_id;
        self.targets.push(target.clone());
        Ok(target)
    }

    /// validate target search structure
    pub fn validate_target<'a>(&self, target: &'a Target, params: &'a TargetSearch) -> bool {
        if params.target_custom_id.is_some() {
            if target.target_custom_id.ne(&params.target_custom_id) {
                return false;
            }
        }

        if let Some(ref name) = params.target_name {
            if target.target_name.ne(name) {
                return false;
            }
        }

        if let Some(ref target_uuid) = params.target_uuid {
            let s_id = target.target_uuid.clone().to_string();
            if s_id.ne(target_uuid) {
                return false;
            }
        };

        if let Some(ref target_id) = params.target_id {
            let s_id = target.target_id.clone().to_string();
            if s_id.ne(target_id) {
                return false;
            }
        };

        if let Some(ref group_uuid) = params.target_group_uuid {
            let search_uuid = uuid::Uuid::from_str(&group_uuid);
            if search_uuid.is_ok() {
                let search_uuid = search_uuid.unwrap();
                if !target.target_groups.contains(&search_uuid) {
                    return false;
                }
            }
        }

        if let Some(ref target_type) = params.target_type {
            if target.target_type.ne(target_type) {
                return false;
            }
        }
        true
    }

    /// checking if target exist
    pub fn exist_target<'a>(&self, params: &'a TargetSearch) -> bool {
        for target in &self.targets {
            if self.validate_target(&target, &params) {
                return true;
            }
        }
        false
    }

    /// delete targets
    pub fn delete_targets(&mut self, params: TargetSearch) {
        let mut targets = self.targets.clone();
        targets.retain(|target| !self.validate_target(target, &params));
        self.targets = targets;
    }

    /// update a target
    pub fn update_targets(&mut self, params: TargetSearch) {
        for mut target in self.targets.iter_mut() {
            if params.target_custom_id.is_some() {
                if target.target_custom_id.ne(&params.target_custom_id) {
                    continue;
                }
            } else {
                if let Some(ref id) = params.target_id {
                    let s_id = target.target_id.clone().to_string();
                    if s_id.ne(id) {
                        continue;
                    }
                } else {
                    continue;
                }
            }

            if let Some(ref name) = params.target_name {
                target.target_name = name.clone();
            }

            if let Some(ref target_type) = params.target_type {
                target.target_type = target_type.clone();
            }

            if let Some(ref meta) = params.meta {
                for (k, v) in meta {
                    target.meta.insert(k.clone(), v.clone());
                }
            }
        }
    }

    /// get targets with params
    pub fn get_targets(&self, params: TargetSearch) -> Vec<Target> {
        let mut results: Vec<Target> = vec![];

        for target in &self.targets {
            if self.validate_target(&target, &params) {
                results.push(target.clone());
            }
        }
        results
    }

    /// get targets by group uuid
    pub fn get_targets_by_group<'a>(&self, group: &'a Group) -> Vec<Target> {
        let mut search = TargetSearch::default();
        search.target_group_uuid = Some(group.group_uuid.to_string());
        self.get_targets(search)
    }
}
