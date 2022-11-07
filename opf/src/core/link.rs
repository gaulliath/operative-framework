use super::session::Session;
use crate::common::search::Link as LinkSearch;
use crate::common::{self, Link};
use std::str::FromStr;

impl Session {
    pub fn create_link(&mut self, link: Link) {
        self.links.push(link);
    }

    pub fn validate_link<'a>(&self, link: &'a Link, params: &'a LinkSearch) -> bool {
        let mut validate = true;

        if let Some(ref target) = params.link_target {
            match uuid::Uuid::from_str(target) {
                Ok(uuid) => {
                    validate = link.link_target.eq(&uuid);
                }
                Err(_) => {
                    return false;
                }
            }
        }

        if let Some(ref source) = params.link_source {
            match uuid::Uuid::from_str(source) {
                Ok(uuid) => {
                    validate = link.link_source.eq(&uuid);
                }
                Err(_) => {
                    return false;
                }
            }
        }

        if let Some(ref label) = params.link_label {
            if link.link_label.ne(label) {
                return false;
            }
        }

        if params.link_color.is_some() {
            if link.link_color.ne(&params.link_color) {
                return false;
            }
        }

        if let Some(ref link_id) = params.link_id {
            if let Ok(uuid) = uuid::Uuid::from_str(link_id) {
                validate = link.link_id.eq(&uuid);
            }
        };

        validate
    }

    pub fn exist_link(&self, params: LinkSearch) -> bool {
        for link in &self.links {
            if self.validate_link(&link, &params) {
                return true;
            }
        }
        false
    }

    pub fn delete_links(&mut self, params: LinkSearch) {
        let mut links = self.links.clone();
        links.retain(|link| !self.validate_link(link, &params));
        self.links = links;
    }

    pub fn update_links(&mut self, params: LinkSearch) {
        for mut link in self.links.iter_mut() {
            if let Some(ref id) = params.link_id {
                match uuid::Uuid::from_str(id.clone().as_str()) {
                    Ok(id) => {
                        if link.link_id.ne(&id) {
                            continue;
                        }
                    }
                    Err(_) => continue,
                };
            }

            if let Some(ref label) = params.link_label {
                link.link_label = label.clone();
            }

            if params.link_color.is_some() {
                link.link_color = params.link_color.clone();
            }
        }
    }

    pub fn get_links(&self, params: LinkSearch) -> Vec<Link> {
        let mut results: Vec<Link> = vec![];

        for link in &self.links {
            if self.validate_link(&link, &params) {
                results.push(link.clone());
            }
        }
        results
    }
}
