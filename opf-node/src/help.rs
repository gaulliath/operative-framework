use opf_models::error::ErrorKind;
use opf_models::TargetType;
use opf_models::event::{Domain, Event, send_event_to};
use strum::IntoEnumIterator;
use crate::Node;

unsafe impl Sync for Node {}
unsafe impl Send for Node {}

impl Node {
    pub async fn on_help(&self) -> Result<(), ErrorKind> {
        let headers = vec![
            "command".to_string(),
            "description".to_string(),
        ];
        let mut rows = vec![];

        for (name, description, usage) in opf_models::command::LISTS {
            let fields = vec![
                name.to_string(),
                description.to_string(),
                usage.to_string()
            ];
            rows.push(fields);
        }
        let _ = send_event_to(&self.self_tx, (Domain::CLI, Event::ResponseTable((headers, rows)))).await;

        // listing of target type
        let headers = vec![
            "target_type available".to_string(),
        ];
        let mut rows = vec![];
        for target_type in TargetType::iter() {
            let el = vec![target_type.to_string()];
            rows.push(el);
        }
        let _ = send_event_to(&self.self_tx, (Domain::CLI, Event::ResponseTable((headers, rows)))).await;
        Ok(())
    }
}