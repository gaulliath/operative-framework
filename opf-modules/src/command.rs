use opf_models::error::ErrorKind;
use opf_models::event::Event::ResponseSimple;
use opf_models::event::{send_event_to, Domain};
use std::collections::HashMap;

use crate::{Module, ModuleController};

impl ModuleController {
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
