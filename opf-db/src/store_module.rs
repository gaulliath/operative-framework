use opf_models::event::{send_event, Event};
use opf_models::{
    error::{ErrorKind, Group as GroupError, Module as ModuleError},
    module, Command, Group, KeyStore, Target,
};
use std::collections::HashMap;

use crate::store::DB;

impl DB {
    pub async fn on_prepare_module(
        &mut self,
        keystore: &KeyStore,
        module_name: String,
        command: Command,
    ) -> Result<(), ErrorKind> {
        let mut groups = self.groups.write().await;
        let targets = self.targets.read().await;
        let mut arguments = command.params;

        let group_id = (groups.len() + 1) as i32;
        let new_group = Group::new(
            group_id,
            format!("{}_{}", module_name, chrono::Utc::now().to_rfc3339()),
        );
        groups.insert(group_id, new_group);

        if let Some(n_group_id) = arguments.get("group_id") {
            let search_group = n_group_id
                .clone()
                .parse::<i32>()
                .map_err(|_| ErrorKind::InvalidCommandArgument)?;

            match groups.get(&search_group) {
                Some(group) => {
                    let target_ids = group.targets.clone();
                    for id in target_ids {
                        if let Some(target) = targets.get(&id) {
                            let mut args = arguments.clone();
                            args.insert(module::TARGET.to_string(), target.target_name.clone());
                            args.insert(
                                module::TARGET_ID.to_string(),
                                target.target_id.to_string(),
                            );
                            self.prepare_data(target, &mut args, keystore).await?;

                            let _ = send_event(
                                &self.db_tx,
                                Event::ExecuteModule((
                                    group_id,
                                    target.clone(),
                                    module_name.clone(),
                                    args,
                                )),
                            )
                            .await;
                        }
                    }

                    return Ok(());
                }
                None => return Err(ErrorKind::InvalidCommandArgument),
            }
        }

        let target = arguments
            .get(module::TARGET_ID)
            .map(|id| id.parse::<i32>())
            .map(|id| targets.get(&id.ok()?))
            .ok_or(ErrorKind::Module(ModuleError::TargetNotAvailable))?
            .ok_or(ErrorKind::Module(ModuleError::TargetNotAvailable))?;

        arguments.insert(module::TARGET.to_string(), target.target_name.clone());
        self.prepare_data(target, &mut arguments, keystore).await?;

        send_event(
            &self.db_tx,
            Event::ExecuteModule((group_id, target.clone(), module_name, arguments)),
        )
        .await
    }

    async fn prepare_data(
        &self,
        target: &Target,
        arguments: &mut HashMap<String, String>,
        keystore: &KeyStore,
    ) -> Result<(), ErrorKind> {
        for (_, value) in arguments {
            if value.contains("meta:") {
                let need = value.replace("meta:", "");
                let metadata_value = target.meta.get(&need).ok_or(ErrorKind::Module(
                    ModuleError::ParamNotAvailable(format!(
                        "metadata '{}' not available in target",
                        need
                    )),
                ))?;
                *value = metadata_value.clone();
            } else if value.contains("keystore:") {
                let need = value.replace("keystore:", "");
                let key_value = keystore.0.get(&need).ok_or(ErrorKind::Module(
                    ModuleError::ParamNotAvailable(format!(
                        "key '{}' not available in keystore",
                        need
                    )),
                ))?;
                *value = key_value.clone();
            }
        }
        Ok(())
    }

    pub async fn on_results_targets(
        &mut self,
        group_id: i32,
        targets: Vec<Target>,
    ) -> Result<(), ErrorKind> {
        let mut groups = self.groups.write().await;
        let current_group = groups
            .get_mut(&group_id)
            .ok_or(ErrorKind::Group(GroupError::Exist(group_id.to_string())))?;

        for mut target in targets {
            if self
                .target_exist(&target.target_name, &target.target_type)
                .await
            {
                continue;
            }
            let mut targets = self.targets.write().await;

            let id = (targets.len() + 1) as i32;
            target.target_id = id;
            targets.insert(id, target);
            current_group.targets.push(id);
        }
        Ok(())
    }
}
