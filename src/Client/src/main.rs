use crate::{app::App, client::CommanderClient, config::Config};

mod action;
mod app;
mod client;
mod components;
mod config;

#[tokio::main]
async fn main() -> Result<(), client::ClientError> {
    let config = Config::from_env();
    let client = Box::new(CommanderClient::connect(&config.commander_addr).await?);

    let mut app = App::new(config, client);
    let mut terminal = ratatui::init();

    let result = app.run(&mut terminal).await;

    ratatui::restore();
    result
}
