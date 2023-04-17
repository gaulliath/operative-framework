use std::collections::HashMap;

use tokio::sync::mpsc::UnboundedSender;

use opf_models::error::ErrorKind;
use opf_models::event::send_event;
use opf_models::event::Event;
use opf_models::metadata::Args;

use crate::CompiledModule;

pub async fn run(
    controller_tx: UnboundedSender<Event>,
    module: Box<dyn CompiledModule>,
    args: HashMap<String, String>,
) -> Result<(), ErrorKind> {
    let mut params = module.args();
    for param in params.iter_mut() {
        match args.get(&param.name) {
            Some(value) => param.value = Some(value.clone()),
            None => {
                if !param.is_optional {
                    return send_event(
                        &controller_tx,
                        Event::ResponseError(format!("argument '{}' not available", param.name)),
                    )
                    .await;
                }
            }
        }
    }

    if module.is_threaded() {
        return match module
            .run(Args::from(params), Some(controller_tx.clone()))
            .await
        {
            Ok(_) => {
                send_event(
                    &controller_tx,
                    Event::ResponseSimple(format!("module threaded executed")),
                )
                .await
            }
            Err(e) => send_event(&controller_tx, Event::ResponseError(e.to_string())).await,
        };
    }

    match module.run(Args::from(params), None).await {
        Ok(targets) => {
            let _ = send_event(
                &controller_tx,
                Event::ResponseSimple(format!(
                    "'{}' targets add from module '{}'",
                    module.name(),
                    targets.len()
                )),
            )
            .await;
            send_event(&controller_tx, Event::ResultsModule(targets)).await
        }
        Err(e) => send_event(&controller_tx, Event::ResponseError(e.to_string())).await,
    }
}
