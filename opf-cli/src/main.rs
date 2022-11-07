mod prompt;
use crossterm::event::{KeyCode, KeyModifiers};
use reedline::ExampleHighlighter;
use reedline::{ColumnarMenu, DefaultCompleter, ReedlineMenu, default_emacs_keybindings, Emacs};
use reedline::{ListMenu, Reedline, ReedlineEvent, Signal};

pub fn main() {
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
                .required(false)
        )
        .get_matches();

    let (session, mut config) = opf::core::session::new();

    if matches.is_present("verbose") {
        config.verbose = true;
    }

    let mut opf = opf::core::new(session, config);
    let module_path = matches.value_of("modules").unwrap_or("~/.opf/modules/");
    opf.init_manager(module_path);

    let auto_complete = vec![
        "add target".into(),
        "list".into(),
        "list modules".into(),
        "list target".into(),
        "list group".into(),
    ];
    let completer = Box::new(DefaultCompleter::new_with_wordlen(auto_complete.clone(), 2));
    // Use the interactive menu to select options from the completer
    let completion_menu = Box::new(
        ColumnarMenu::default()
            .with_name("completion_menu")
            .with_only_buffer_difference(false),
    );

    let mut line_editor = Reedline::create()
        .with_completer(completer)
        .with_quick_completions(true)
        .with_partial_completions(true)
        .with_menu(ReedlineMenu::EngineCompleter(completion_menu))
        .with_edit_mode({
            let mut keybindings = default_emacs_keybindings();
            keybindings.add_binding(
                reedline::KeyModifiers::NONE,
                reedline::KeyCode::Tab,
                ReedlineEvent::UntilFound(vec![
                    ReedlineEvent::Menu("completion_menu".to_string()),
                    ReedlineEvent::MenuNext,
                ]),
            );

            Box::new(Emacs::new(keybindings))
        });
    //.with_highlighter(Box::new(ExampleHighlighter::new(commands)));
    let prompt = prompt::PromptCli {};
    loop {
        let sig = line_editor.read_line(&prompt).unwrap();
        match sig {
            Signal::Success(buffer) => {
                if let Err(e) = opf.run(buffer.trim().to_string()) {
                    println!("ERR {}", e);
                }
            }
            Signal::CtrlD | Signal::CtrlC => {
                return;
            }
        }
    }
}
