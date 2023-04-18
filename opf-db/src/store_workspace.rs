use std::time::SystemTime;

use opf_models::event::{send_event, send_event_to, Domain, Event};
use opf_models::{
    self,
    error::{ErrorKind, Workspace as WorkspaceError},
    workspace, Command, CommandAction, Workspace,
};

use crate::store::DB;
use crate::DBStore;

impl DBStore {
    pub async fn on_workspace_command(&mut self, command: Command) -> Result<(), ErrorKind> {
        match command.action {
            CommandAction::Add => self.add_workspace(command).await,
            CommandAction::Switch => self.switch_workspace(command).await,
            CommandAction::List => self.list_workspaces(command).await,
            //CommandAction::Del => self.delete_target(command).await,
            CommandAction::Set => self.update_workspace(command).await,
            _ => Err(ErrorKind::ActionNotAvailable),
        }
    }

    async fn add_workspace(&mut self, cmd: Command) -> Result<(), ErrorKind> {
        let params = cmd.params;
        let workspaces = &mut self.workspaces;

        let workspace_id = workspaces.len() as i32;
        let workspace_name = params
            .get(workspace::NAME)
            .ok_or(ErrorKind::Workspace(WorkspaceError::ParamNotFound(
                workspace::NAME.to_string(),
            )))?
            .clone();

        for (_, workspace) in workspaces.iter() {
            if workspace.workspace_name.eq(&workspace_name) {
                return Err(ErrorKind::Workspace(WorkspaceError::Exist(
                    workspace_name.to_string(),
                )));
            }
        }

        let new_workspace = Workspace {
            workspace_id,
            workspace_name: workspace_name.clone(),
            workspace_created_at: SystemTime::now(),
        };

        workspaces.insert(workspace_id, new_workspace.clone());
        self.dbs.insert(workspace_id, DB::new(self.self_tx.clone()));

        send_event(
            &self.self_tx,
            Event::ResponseSimple(format!("workspace created => {}", workspace_name)),
        )
        .await
    }

    async fn switch_workspace(&mut self, cmd: Command) -> Result<(), ErrorKind> {
        let workspace_id = cmd
            .params
            .get(workspace::ID)
            .ok_or(ErrorKind::InvalidFormatArgument)?;

        let workspace_id = workspace_id
            .parse::<i32>()
            .map_err(|_| ErrorKind::InvalidFormatArgument)?;

        let workspace = self
            .workspaces
            .get(&workspace_id)
            .ok_or(ErrorKind::Workspace(WorkspaceError::GenericError(format!(
                "workspace id '{}' not exist",
                workspace_id
            ))))?;

        self.current_workspace = workspace_id;

        send_event_to(
            &self.node_tx,
            (
                Domain::CLI,
                Event::SetWorkspace(workspace.workspace_name.clone()),
            ),
        )
        .await
    }

    async fn update_workspace(&mut self, cmd: Command) -> Result<(), ErrorKind> {
        let workspace_id = cmd
            .params
            .get(workspace::ID)
            .ok_or(ErrorKind::InvalidFormatArgument)?;

        let workspace_id = workspace_id
            .parse::<i32>()
            .map_err(|_| ErrorKind::InvalidFormatArgument)?;

        let mut workspace = self
            .workspaces
            .get_mut(&workspace_id)
            .ok_or(ErrorKind::Workspace(WorkspaceError::GenericError(format!(
                "workspace not found '{}'",
                workspace_id
            ))))?;

        for (field, value) in cmd.params {
            match field.as_str() {
                workspace::NAME => workspace.workspace_name = value,
                _ => {}
            }
        }

        let _ = send_event_to(
            &self.node_tx,
            (
                Domain::CLI,
                Event::SetWorkspace(workspace.workspace_name.clone()),
            ),
        )
        .await;

        send_event(
            &self.self_tx,
            Event::ResponseSimple(format!(
                "workspace updated successfully : {:#?}",
                &workspace.workspace_name
            )),
        )
        .await
    }

    async fn list_workspaces(&self, _cmd: Command) -> Result<(), ErrorKind> {
        let workspaces = &self.workspaces;

        let mut headers = vec![
            workspace::ID.to_string(),
            workspace::NAME.to_string(),
            workspace::CREATED_AT.to_string(),
        ];

        let mut header_meta = vec![];
        let mut rows = vec![];

        for (_, workspace) in workspaces.iter() {
            let created_at: chrono::DateTime<chrono::Utc> = workspace.workspace_created_at.into();
            let fields = vec![
                workspace.workspace_id.to_string(),
                workspace.workspace_name.to_string(),
                created_at.to_rfc3339(),
            ];
            rows.push(fields);
        }

        headers.append(&mut header_meta);
        send_event(&self.self_tx, Event::ResponseTable((headers, rows))).await
    }
}
