use opf_models::event::{send_event, Event};
use opf_models::{self, error::ErrorKind, group, Command, CommandAction};

use crate::store::DB;

impl DB {
    pub async fn on_group_command(&mut self, command: Command) -> Result<(), ErrorKind> {
        match command.action {
            CommandAction::List => self.list_groups(command).await,
            _ => Err(ErrorKind::ActionNotAvailable),
        }
    }

    async fn list_groups(&self, _cmd: Command) -> Result<(), ErrorKind> {
        let groups = &self.groups.read().await;

        let mut headers = vec![
            group::ID.to_string(),
            group::NAME.to_string(),
            group::TARGETS.to_string(),
        ];
        let mut header_meta = vec![];
        let mut rows = vec![];

        for (_, group) in groups.iter() {
            let fields = vec![
                group.group_id.to_string(),
                group.group_name.to_string(),
                group.targets.len().to_string(),
            ];
            rows.push(fields);
        }

        headers.append(&mut header_meta);
        send_event(&self.db_tx, Event::ResponseTable((headers, rows))).await
    }
}
