using Commander.Core.Factories;
using Commander.Core.Ports;
using Commander.Infrastructure.Adapters;
using Commander.Server.Services;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddGrpc();
builder.Services.AddGrpcReflection();

builder.Services.AddSingleton<IRunnerPort, DockerRunnerAdapter>();
builder.Services.AddSingleton<JobDefinitionFactory>();

var app = builder.Build();

if (app.Environment.IsDevelopment())
{
  app.MapGrpcReflectionService();
}

app.MapGrpcService<OrchestratorService>();
app.MapGet("/", () => "Communication with gRPC endpoints must be made through a gRPC client.");

app.Run();
