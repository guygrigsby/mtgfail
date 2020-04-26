use actix_web::{get, web, App, HttpRequest, HttpResponse, HttpServer, Responder};
use scryfall::card::Card;

#[actix_rt::main]
async fn main() -> std::io::Result<()> {
    HttpServer::new(|| {
        App::new()
            .route("/", web::get().to(greet))
            .route("/{name}", web::get().to(card_comp))
    })
    .bind("127.0.0.1:8080")?
    .run()
    .await
}
async fn card_comp(req: HttpRequest) -> impl Responder {
    let card_name = req.match_info().get("name").unwrap_or("World");
    match Card::named_fuzzy(card_name) {
        Ok(card) => web::Json(card),
        Err(e) => panic!(format!("{:?}", e)),
    }
}
#[get("/")]
async fn index3() -> impl Responder {
    HttpResponse::Ok().body("mtg.fail")
}
async fn greet(req: HttpRequest) -> impl Responder {
    let name = req.match_info().get("name").unwrap_or("World");
    format!("Hello {}!", &name)
}
