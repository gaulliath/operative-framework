use opf_models::event::{send_event, Event};
use opf_models::{
    error::{ErrorKind, Module as ModuleError},
    module, Command, Target,
};

use crate::store::DB;

impl DB {
    pub async fn on_prepare_module(
        &mut self,
        module_name: String,
        command: Command,
    ) -> Result<(), ErrorKind> {
        let targets = self.targets.read().await;
        let mut arguments = command.params;

        let target = arguments
            .get(module::TARGET_ID)
            .map(|id| id.parse::<i32>())
            .map(|id| targets.get(&id.ok()?))
            .ok_or(ErrorKind::Module(ModuleError::TargetNotAvailable))?
            .ok_or(ErrorKind::Module(ModuleError::TargetNotAvailable))?;

        arguments.insert(module::TARGET.to_string(), target.target_name.clone());

        for (_, value) in &mut arguments {
            if value.contains("meta:") {
                let need = value.replace("meta:", "");
                let metadata_value = target.meta.get(&need).ok_or(ErrorKind::Module(
                    ModuleError::ParamNotAvailable(format!(
                        "metadata '{}' not available in target",
                        need
                    )),
                ))?;
                *value = metadata_value.clone();
            }
        }

        send_event(&self.db_tx, Event::ExecuteModule((module_name, arguments))).await
    }

    pub async fn on_results_targets(&mut self, targets: Vec<Target>) -> Result<(), ErrorKind> {
        for mut target in targets {
            if self
                .target_exist(&target.target_name, &target.target_type)
                .await
            {
                continue;
            }
            let mut targets = self.targets.write().await;

            let id = targets.len() + 1;
            target.target_id = id as i32;
            targets.insert(id as i32, target);
        }
        Ok(())
    }
}
