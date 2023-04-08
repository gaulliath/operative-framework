use colored::Colorize;
use rustyline::completion::FilenameCompleter;
use rustyline::error::ReadlineError;
use rustyline::highlight::{Highlighter, MatchingBracketHighlighter};
use rustyline::hint::HistoryHinter;
use rustyline::validate::MatchingBracketValidator;
use rustyline::{Cmd, CompletionType, Config, EditMode, Editor, KeyEvent};
use rustyline::{Completer, Helper, Hinter, Result, Validator};
use std::borrow::Cow::{self, Borrowed, Owned};

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

pub fn main() -> Result<()> {
    let matches = clap::App::new("operative-framework")
        .author("Tristan Granier")
        .about("OSINT framework with LUA interpreter.")
        .arg(
            clap::Arg::new("modules")
                .required(false)
                .long("modules")
                .takes_value(true)
                .help("location path of LUA modules."),
        )
        .arg(
            clap::Arg::new("verbose")
                .short('v')
                .help("print logging verbose")
                .required(false)
                .takes_value(false),
        )
        .arg(
            clap::Arg::new("cli")
                .help("use with interactive shell")
                .long("cli")
                .required(false),
        )
        .get_matches();

    let (session, mut config) = opf::core::session::new();

    if matches.is_present("verbose") {
        config.verbose = true;
    }

    let mut opf = opf::core::new(session, config);
    let module_path = matches.value_of("modules").unwrap_or("~/.opf/modules/");
    if let Err(e) = opf.init_manager(module_path) {
        println!("ERR {}", e);
    }

    let config = Config::builder()
        .history_ignore_space(true)
        .completion_type(CompletionType::List)
        .edit_mode(EditMode::Emacs)
        .build();
    let h = MyHelper {
        completer: FilenameCompleter::new(),
        highlighter: MatchingBracketHighlighter::new(),
        hinter: HistoryHinter {},
        colored_prompt: format!("{}", "(opf) > ".yellow()).to_owned(),
        validator: MatchingBracketValidator::new(),
    };
    let mut rl = Editor::with_config(config)?;
    rl.set_helper(Some(h));
    rl.bind_sequence(KeyEvent::alt('n'), Cmd::HistorySearchForward);
    rl.bind_sequence(KeyEvent::alt('p'), Cmd::HistorySearchBackward);
    if rl.load_history("/tmp/history.txt").is_err() {
        println!("No previous history.");
    }

    loop {
        let readline = rl.readline("(opf) > ");
        match readline {
            Ok(line) => {
                rl.add_history_entry(line.as_str())?;
                if let Err(e) = opf.run(&line) {
                    println!("ERR {}", e);
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
