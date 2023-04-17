use std::borrow::Cow::{self, Borrowed, Owned};
use std::sync::{Arc, RwLock};
use std::time::Duration;

use ::log::{error, info};
use colored::Colorize;
use comfy_table::Table;
use rustyline::completion::FilenameCompleter;
use rustyline::error::ReadlineError;
use rustyline::highlight::{Highlighter, MatchingBracketHighlighter};
use rustyline::hint::HistoryHinter;
use rustyline::validate::MatchingBracketValidator;
use rustyline::{Completer, Helper, Hinter, Result, Validator};
use rustyline::{CompletionType, Config, EditMode, Editor, ExternalPrinter};

use opf_models::event::{send_event_to, Domain, Event};
use opf_node::Node;

mod config;
mod log;

//use crossterm::style::Stylize;

#[derive(Helper, Completer, Hinter, Validator)]
struct MyHelper {
    #[rustyline(Completer)]
    completer: FilenameCompleter,
    highlighter: MatchingBracketHighlighter,
    #[rustyline(Validator)]
    validator: MatchingBracketValidator,
    #[rustyline(Hinter)]
    hinter: HistoryHinter,
    colored_prompt: String,
}

impl Highlighter for MyHelper {
    fn highlight_prompt<'b, 's: 'b, 'p: 'b>(
        &'s self,
        prompt: &'p str,
        default: bool,
    ) -> Cow<'b, str> {
        if default {
            Borrowed(&self.colored_prompt)
        } else {
            Borrowed(prompt)
        }
    }

    fn highlight_hint<'h>(&self, hint: &'h str) -> Cow<'h, str> {
        Owned("\x1b[1m".to_owned() + hint + "\x1b[m")
    }

    fn highlight<'l>(&self, line: &'l str, pos: usize) -> Cow<'l, str> {
        self.highlighter.highlight(line, pos)
    }

    fn highlight_char(&self, line: &str, pos: usize) -> bool {
        self.highlighter.highlight_char(line, pos)
    }
}

#[tokio::main]
async fn main() -> Result<()> {
    let _ = log::init();
    let (tx, rx) = std::sync::mpsc::channel::<Event>();
    let (node_tx, node) = Node::new(tx).await;
    let workspace = Arc::new(RwLock::new(String::from("test-123")));

    let config = Config::builder()
        .history_ignore_space(true)
        .completion_type(CompletionType::List)
        .edit_mode(EditMode::Emacs)
        .build();
    let h = MyHelper {
        completer: FilenameCompleter::new(),
        highlighter: MatchingBracketHighlighter::new(),
        hinter: HistoryHinter {},
        colored_prompt: "".to_string(),
        validator: MatchingBracketValidator::new(),
    };
    let mut rl = Editor::with_config(config)?;
    rl.set_helper(Some(h));
    if rl.load_history("/tmp/history.txt").is_err() {
        info!("history file '/tmp/history.txt' not found, we create one.");
    }

    let mut external_printer = rl
        .create_external_printer()
        .expect("cant't create external printer");

    let workspace_clone = Arc::clone(&workspace);
    std::thread::spawn(move || loop {
        for event in rx.iter() {
            match event {
                Event::ResponseSimple(res) => {
                    let _ = external_printer.print(format!("INF {}", res));
                }
                Event::ResponseError(res) => {
                    let _ = external_printer.print(format!("ERR {}", res));
                }
                Event::ResponseTable((headers, rows)) => {
                    let mut table = Table::new();
                    table.set_header(headers);
                    table.add_rows(rows);
                    let _ = external_printer.print(format!("{}", table));
                }
                Event::SetWorkspace(new_workspace) => {
                    let mut workspace = workspace_clone.write().unwrap();
                    *workspace = new_workspace.clone();
                    let _ = external_printer.print(format!("workspace => {}", new_workspace));
                }
                _ => {}
            }
        }
    });

    tokio::spawn(node.main_loop());

    std::thread::sleep(Duration::from_secs(1));

    loop {
        let workspace_mame = { workspace.read().unwrap().clone() };
        let prompt = format!(
            "[{} [workspace:{}]> ",
            "opf".yellow(),
            workspace_mame.blue()
        );
        rl.helper_mut().unwrap().colored_prompt = prompt.clone();
        let readline = rl.readline(prompt.as_str());

        match readline {
            Ok(line) => {
                rl.add_history_entry(line.as_str())?;
                if let Err(e) =
                    send_event_to(&node_tx, (Domain::Node, Event::NewCommand(line))).await
                {
                    error!("{}", e);
                }
            }
            Err(ReadlineError::Interrupted) => {
                println!("Interrupted");
                break;
            }
            Err(ReadlineError::Eof) => {
                println!("Encountered Eof");
                break;
            }
            Err(err) => {
                println!("Error: {err:?}");
                break;
            }
        }
    }

    rl.append_history("/tmp/history.txt")
}
