pub fn init() -> Result<(), log::SetLoggerError> {
    fern::Dispatch::new()
        // Perform allocation-free log formatting
        .format(|out, message, record| {
            out.finish(format_args!(
                "[{} {}] {}",
                humantime::format_rfc3339(std::time::SystemTime::now()),
                record.level(),
                message
            ))
        })
        // Add blanket level filter -
        .level(log::LevelFilter::Info)
        // - and per-module overrides
        .level_for("hyper", log::LevelFilter::Info)
        // Output to stdout, files, and other Dispatch configurations
        .chain(std::io::stdout())
        .chain(fern::log_file("output.log").unwrap())
        // Apply globally
        .apply()
        .expect("cant't be parse error");

    Ok(())
}
