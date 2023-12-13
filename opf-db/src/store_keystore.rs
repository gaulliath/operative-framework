use crate::DBStore;
use opf_models::event::{send_event, Event};
use opf_models::{self, error::ErrorKind, keystore, Command, CommandAction};

impl DBStore {
    pub async fn on_keystore_command(&mut self, command: Command) -> Result<(), ErrorKind> {
        match command.action {
            CommandAction::List => self.list_keystore(command).await,
            _ => Err(ErrorKind::ActionNotAvailable),
        }
    }

    async fn list_keystore(&self, _cmd: Command) -> Result<(), ErrorKind> {
        let keystore = &self.keystore.read().await;

        let mut headers = vec![keystore::NAME.to_string(), keystore::VALUE.to_string()];

        let mut header_meta = vec![];
        let mut rows = vec![];

        for (name, value) in &keystore.0 {
            let fields = vec![name.clone(), value.clone()];
            rows.push(fields);
        }

        headers.append(&mut header_meta);
        send_event(&self.self_tx, Event::ResponseTable((headers, rows))).await
    }
}
