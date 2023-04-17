use opf_models::event::{send_event, Event};
use opf_models::{
    self,
    error::{ErrorKind, Target as TargetError},
    target, Command, CommandAction, Target, TargetType,
};

use crate::store::DB;

impl DB {
    pub async fn on_target_command(&mut self, command: Command) -> Result<(), ErrorKind> {
        match command.action {
            CommandAction::Add => self.add_target(command).await,
            CommandAction::List => self.list_targets(command).await,
            CommandAction::Del => self.delete_target(command).await,
            CommandAction::Set => self.update_target(command).await,
            _ => Err(ErrorKind::ActionNotAvailable),
        }
    }

    pub async fn target_exist(&self, name: &str, target_type: &TargetType) -> bool {
        let targets = self.targets.read().await;
        for (_, target) in targets.iter() {
            if target.target_name.eq(name) && target.target_type.eq(target_type) {
                return true;
            }
        }
        false
    }

    async fn update_target(&mut self, cmd: Command) -> Result<(), ErrorKind> {
        let target_id = cmd
            .params
            .get(target::ID)
            .ok_or(ErrorKind::InvalidFormatArgument)?;

        let target_id = target_id
            .parse::<i32>()
            .map_err(|_| ErrorKind::InvalidFormatArgument)?;

        let mut targets = self.targets.write().await;

        let target = targets
            .get_mut(&target_id)
            .ok_or(ErrorKind::Target(TargetError::NotFound))?;

        let mut data = target.to_hashmap();
        let mut updated_data = cmd.params;

        updated_data.remove(target::ID);

        data.extend(updated_data);
        *target = Target::try_from(data).map_err(|_| ErrorKind::InvalidFormatArgument)?;

        send_event(
            &self.db_tx,
            Event::ResponseSimple(format!("target updated successfully : {:#?}", target)),
        )
        .await
    }

    async fn add_target(&mut self, cmd: Command) -> Result<(), ErrorKind> {
        let mut params = cmd.params;
        let mut targets = self.targets.write().await;

        let target_name = match params.get(target::NAME) {
            Some(name) => name.clone(),
            None => return Err(ErrorKind::Target(TargetError::ParamNameNotFound)),
        };

        let target_type = match params.get(target::TYPE) {
            Some(t) => opf_models::validate_type(t)?,
            None => return Err(ErrorKind::Target(TargetError::ParamTypeNotFound)),
        };

        params.remove(target::NAME);
        params.remove(target::TYPE);

        let target_parent = params
            .get(target::PARENT)
            .map(|id| id.parse::<i32>().unwrap_or(0));

        let target_custom_id = params.get("custom_id").map(|custom| custom.to_string());

        for (_, exist_target) in targets.iter() {
            if exist_target.target_name.eq(&target_name)
                && exist_target.target_type.eq(&target_type)
            {
                return Err(ErrorKind::Target(TargetError::TargetExist));
            }
        }

        let target_id = targets.len() as i32 + 1;

        let new_target = Target {
            target_id: target_id.clone(),
            target_name,
            target_type,
            target_custom_id,
            target_groups: vec![],
            target_parent,
            meta: params.clone(),
        };

        targets.insert(target_id, new_target.clone());

        send_event(
            &self.db_tx,
            Event::ResponseSimple(format!("{:#?}", new_target)),
        )
        .await
    }

    async fn delete_target(&mut self, cmd: Command) -> Result<(), ErrorKind> {
        let target_id = cmd
            .params
            .get(target::ID)
            .ok_or(ErrorKind::InvalidFormatArgument)
            .map(|id| {
                id.parse::<i32>()
                    .map_err(|_| ErrorKind::InvalidFormatArgument)
            })??;

        let mut message =
            Event::ResponseSimple(format!("target with id={} successfully deleted", target_id));
        if self.targets.write().await.remove(&target_id).is_none() {
            message = Event::ResponseError(format!("target with id={} not found", target_id));
        }
        send_event(&self.db_tx, message).await
    }

    async fn list_targets(&self, cmd: Command) -> Result<(), ErrorKind> {
        let params = cmd.params;
        let targets = &self.targets.read().await;

        let with_metadata = params.get("show").is_some();

        let mut headers = vec![
            target::ID.to_string(),
            target::PARENT.to_string(),
            target::CUSTOM_ID.to_string(),
            target::TYPE.to_string(),
            target::NAME.to_string(),
        ];
        let mut header_meta = vec![];
        let mut rows = vec![];

        for (_, target) in targets.iter() {
            let parent = match target.target_parent {
                Some(id) => id.to_string(),
                None => String::from("-"),
            };

            let mut fields = vec![
                target.target_id.to_string(),
                parent,
                target.target_custom_id.clone().unwrap_or(String::from("-")),
                target.target_type.to_string(),
                target.target_name.clone(),
            ];

            if with_metadata {
                for (key, _) in &target.meta.clone() {
                    let key_s = String::from(key);
                    if header_meta.contains(&key_s) {
                        continue;
                    }
                    header_meta.push(key_s);
                }

                for head in &header_meta {
                    match target.meta.get(head) {
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
