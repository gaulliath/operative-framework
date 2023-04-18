use opf_models::error::{ErrorKind, Module as ModuleError};
use opf_models::event::Event::ResponseSimple;
use opf_models::event::{send_event, send_event_to, Domain, Event};
use opf_models::module;
use std::collections::HashMap;

use crate::{Module, ModuleController};

impl ModuleController {
    pub async fn on_help_module(&mut self, module_name: String) -> Result<(), ErrorKind> {
        for (_, module) in &self.modules {
            if module.name().eq(&module_name) {
                let headers = vec![
                    module::NAME.to_string(),
                    module::RESUME.to_string(),
                    module::AUTHOR.to_string(),
                    module::TARGET_TYPE.to_string(),
                ];
                let mut rows = vec![];

                let fields = vec![module.name(), module.author(), module.resume(), module.target_type().to_string()];
                rows.push(fields);
                let _ = send_event(&self.self_tx, Event::ResponseTable((headers, rows))).await;

                let headers = vec![
                    module::ARGUMENT.to_string(),
                    module::ARGUMENT_DEFAULT.to_string(),
                    module::ARGUMENT_OPTIONAL.to_string(),
                ];

                let rows = &module
                    .args()
                    .iter()
                    .filter(|arg| arg.name.ne(module::TARGET))
                    .map(|arg| {
                        vec![
                            arg.name.clone(),
                            arg.value.clone().unwrap_or("<None>".to_string()),
                            format!("{}", arg.is_optional),
                        ]
                    })
                    .collect::<Vec<Vec<String>>>();
                return send_event(&self.self_tx, Event::ResponseTable((headers, rows.clone())))
                    .await;
            }
        }
        Err(ErrorKind::Module(ModuleError::ModuleNameNotFound))
    }

    pub async fn on_list_modules(&mut self) -> Result<(), ErrorKind> {
        let modules = &self.modules;
        let mut headers = vec![
            module::NAME.to_string(),
            module::RESUME.to_string(),
            module::TARGET_TYPE.to_string(),
            module::ARGUMENTS.to_string(),
        ];
        let mut header_meta = vec![];
        let mut rows = vec![];

        for (_, module) in modules.iter() {
            let params = module
                .args()
                .iter()
                .filter(|arg| arg.name.ne(module::TARGET))
                .map(|arg| format!("{}", arg.name))
                .collect::<Vec<String>>();

            let fields = vec![
                module.name(),
                module.resume(),
                module.target_type().to_string(),
                params.join(" | "),
            ];

            rows.push(fields);
        }

        headers.append(&mut header_meta);
        send_event(&self.self_tx, Event::ResponseTable((headers, rows))).await
    }

    pub async fn on_execute_module(
        &mut self,
        data: (String, HashMap<String, String>),
    ) -> Result<(), ErrorKind> {
        for (_, module) in self.modules.iter() {
            if module.name().eq(&data.0) {
                if let Module::Compiled(compiled) = module {
                    let _ = send_event_to(
                        &self.node_tx,
                        (
                            Domain::CLI,
                            ResponseSimple(format!("executing module : {}", module.name())),
                        ),
                    )
                    .await;
                    tokio::spawn(crate::worker::run(
                        self.self_tx.clone(),
                        dyn_clone::clone_box(&**compiled),
                        data.1.clone(),
                    ));
                }
            }
        }
        Ok(())
    }
}
