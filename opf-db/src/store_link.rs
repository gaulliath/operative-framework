use std::time::SystemTime;

use opf_models::event::{send_event, Event};
use opf_models::{
    self,
    error::{ErrorKind, Link as LinkError},
    link, Command, CommandAction, Link, LinkFrom, LinkType,
};

use crate::store::DB;

impl DB {
    pub async fn on_link_command(&mut self, command: Command) -> Result<(), ErrorKind> {
        match command.action {
            CommandAction::Add => self.add_link(command).await,
            CommandAction::List => self.list_links(command).await,
            CommandAction::Del => self.delete_link(command).await,
            CommandAction::Set => self.update_link(command).await,
            _ => Err(ErrorKind::ActionNotAvailable),
        }
    }

    async fn update_link(&mut self, cmd: Command) -> Result<(), ErrorKind> {
        let link_id = cmd
            .params
            .get(link::ID)
            .ok_or(ErrorKind::InvalidFormatArgument)
            .map(|id| {
                id.parse::<i32>()
                    .map_err(|_| ErrorKind::InvalidFormatArgument)
            })??;

        let mut links = self.links.write().await;

        let link = links
            .get_mut(&link_id)
            .ok_or(ErrorKind::Link(LinkError::NotFound))?;

        let mut data = link.to_hashmap();
        let mut updated_data = cmd.params;

        updated_data.remove(link::ID);

        data.extend(updated_data);
        *link = Link::try_from(data).map_err(|_| ErrorKind::InvalidFormatArgument)?;

        send_event(
            &self.db_tx,
            Event::ResponseSimple(format!("link updated successfully : {:#?}", link)),
        )
        .await
    }

    async fn add_link(&mut self, cmd: Command) -> Result<(), ErrorKind> {
        let mut params = cmd.params;
        let mut links = self.links.write().await;
        let targets = self.targets.read().await;

        let link_from = params
            .get(link::FROM)
            .ok_or(ErrorKind::Link(LinkError::ParamFromNotFound))
            .map(|from| {
                from.parse::<i32>()
                    .map_err(|_| ErrorKind::InvalidFormatArgument)
            })??;

        let link_to = params
            .get(link::TO)
            .ok_or(ErrorKind::Link(LinkError::ParamToNotFound))
            .map(|to| {
                to.parse::<i32>()
                    .map_err(|_| ErrorKind::InvalidFormatArgument)
            })??;

        let link_type = params
            .get(link::TYPE)
            .ok_or(ErrorKind::Link(LinkError::ParamTypeNotFound))
            .map(|target_type| opf_models::validate_link_type(target_type))??;

        if link_to == link_from {
            return Err(ErrorKind::Link(LinkError::ParamFormatInvalid(format!(
                "{} / {} is same",
                link::FROM,
                link::TO
            ))));
        }

        if !targets.contains_key(&link_to) || !targets.contains_key(&link_from) {
            return Err(ErrorKind::Link(LinkError::TargetNotFound));
        }

        params.remove(link::FROM);
        params.remove(link::TO);
        params.remove(link::TYPE);

        for (_, link) in links.iter() {
            if link.link_to.eq(&link_to) && link.link_from.eq(&link_from) {
                return Err(ErrorKind::Link(LinkError::LinkExist));
            }
        }

        let link_identifier = match &link_type {
            LinkType::In => vec![link_to.clone()],
            LinkType::Out => vec![link_from.clone()],
            LinkType::Both => vec![link_to.clone(), link_from.clone()],
        };

        let mut inserted = vec![];
        for identifier in link_identifier {
            let link_id = (links.len() + 1) as i32;
            let link = Link {
                link_id,
                link_type: link_type.clone(),
                link_meta: Default::default(),
                link_from,
                link_to,
                link_created_by: LinkFrom::CLI,
                link_created_at: SystemTime::now(),
            };

            links.insert(identifier, link.clone());
            inserted.push(link);
        }

        send_event(
            &self.db_tx,
            Event::ResponseSimple(format!("{:#?}\ntotal = {}", inserted, inserted.len())),
        )
        .await
    }

    async fn delete_link(&mut self, cmd: Command) -> Result<(), ErrorKind> {
        let link_id = cmd
            .params
            .get(link::ID)
            .ok_or(ErrorKind::InvalidFormatArgument)
            .map(|id| {
                id.parse::<i32>()
                    .map_err(|_| ErrorKind::InvalidFormatArgument)
            })??;

        let mut response =
            Event::ResponseSimple(format!("link with id={} successfully deleted", link_id));
        if self.links.write().await.remove(&link_id).is_none() {
            response = Event::ResponseError(format!("link with id={} not found", link_id))
        }
        send_event(&self.db_tx, response).await
    }

    async fn list_links(&self, cmd: Command) -> Result<(), ErrorKind> {
        let params = cmd.params;
        let links = &self.links.read().await;
        let targets = &self.targets.read().await;

        let with_metadata = params.get("show").is_some();

        let mut headers = vec![
            link::ID.to_string(),
            link::TO.to_string(),
            link::FROM.to_string(),
            link::TYPE.to_string(),
            link::CREATED_AT.to_string(),
            link::CREATED_BY.to_string(),
        ];
        let mut header_meta = vec![];
        let mut rows = vec![];

        for (_, link) in links.iter() {
            let created_at: chrono::DateTime<chrono::Utc> = link.link_created_at.into();
            let (target_to, target_from) = {
                (
                    targets
                        .get(&link.link_to)
                        .ok_or(ErrorKind::Link(LinkError::TargetNotFound))
                        .map(|target| target.target_name.clone())?,
                    targets
                        .get(&link.link_from)
                        .ok_or(ErrorKind::Link(LinkError::TargetNotFound))
                        .map(|target| target.target_name.clone())?,
                )
            };

            let mut fields = vec![
                link.link_id.to_string(),
                format!("{} ({})", link.link_to.to_string(), target_to),
                format!("{} ({})", link.link_from.to_string(), target_from),
                link.link_type.to_string(),
                created_at.to_rfc3339(),
                link.link_created_by.to_string(),
            ];

            if with_metadata {
                for (key, _) in &link.link_meta.clone() {
                    let key_s = String::from(key);
                    if header_meta.contains(&key_s) {
                        continue;
                    }
                    header_meta.push(key_s);
                }

                for head in &header_meta {
                    match link.link_meta.get(head) {
                        Some(meta) => {
                            fields.push(meta.clone());
                        }
                        None => fields.push("-".to_string()),
                    }
                }
            }
            rows.push(fields);
        }

        headers.append(&mut header_meta);
        send_event(&self.db_tx, Event::ResponseTable((headers, rows))).await
    }
}
