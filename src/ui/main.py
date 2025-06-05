from rich import print
from rich.console import Console
from rich.progress import Progress, SpinnerColumn, TextColumn, BarColumn, TimeElapsedColumn
from subprocess import run, PIPE
from tkinter import Tk, filedialog
from task_trigger import trigger_task
from dotenv import load_dotenv
import os
import sys
import shlex
import threading
import time
from google.cloud import pubsub_v1

load_dotenv()

PROJECT_ID = os.getenv("PROJECT_ID")
PUB_TOPIC = os.getenv("PUB_TOPIC")
console = Console()


def run_command(cmd):
    result = run(shlex.split(cmd), stdout=PIPE, stderr=PIPE, text=True)
    if result.returncode != 0:
        console.print(f"[red]Error:[/red] {result.stderr.strip()}")
        sys.exit(1)
    return result.stdout.strip()

def select_file():
    root = Tk()
    root.withdraw()
    filepath = filedialog.askopenfilename(title="Select a file to upload")
    if not filepath:
        console.print("[red]No file selected. Exiting.[/red]")
        sys.exit(1)
    return filepath

def upload_with_progress(filepath):
    with Progress(
        SpinnerColumn(),
        TextColumn("[progress.description]{task.description}"),
        BarColumn(),
        TimeElapsedColumn(),
        transient=True,
    ) as progress:
        task = progress.add_task("Uploading to GCS...", total=None)
        upload_cmd = f"python3 gcs-sign.py {filepath} gcp upload"
        upload_output = run_command(upload_cmd)
        progress.update(task, advance=100)
        progress.stop()
    console.print(f"[green]Upload complete.[/green]")

def subscribe_once(file_id):
    subscriber = pubsub_v1.SubscriberClient()
    sub_path = subscriber.subscription_path(PROJECT_ID, f"sub-{file_id}")
    topic_path = f"projects/{PROJECT_ID}/topics/{PUB_TOPIC}"

    console.print(f"[blue]Creating subscription: {sub_path}[/blue]")
    try:
        subscriber.create_subscription(name=sub_path, topic=topic_path)
    except Exception as e:
        console.print(f"[red]Subscription creation failed: {e}[/red]")
        return

    stop_event = threading.Event()

    def callback(message):
        payload = message.data.decode("utf-8")
        console.print(f"[magenta]Received message:[/magenta] {payload}")
        if payload == file_id:
            console.print(f"[bold green]Job completed for file_id:[/bold green] {file_id}")
            message.ack()
            stop_event.set()
            console.print(f"[blue]Deleting subscription: {sub_path}[/blue]")
            subscriber.delete_subscription(subscription=sub_path)

    streaming_pull_future = subscriber.subscribe(sub_path, callback=callback)
    console.print("[green]Listening for job completion message...[/green]")

    try:
        while not stop_event.is_set():
            time.sleep(1)
    except KeyboardInterrupt:
        streaming_pull_future.cancel()
        console.print("[yellow]Stopped listening manually.[/yellow]")


def main():
    console.print("[bold cyan]GCS Upload and Signed URL Generator[/bold cyan]")

    filepath = select_file()

    if not os.path.isfile(filepath):
        console.print(f"[red]File does not exist:[/red] {filepath}")
        sys.exit(1)

    upload_with_progress(filepath)

    filename = os.path.basename(filepath)
    console.print(f"[yellow]Fetching signed download URL for {filename}...[/yellow]")
    download_cmd = f"python3 gcs-sign.py {filename} gcp download"
    download_url = run_command(download_cmd)

    console.print("[cyan]Triggering processing task...[/cyan]")
    file_id = trigger_task(download_url)
    if file_id:
        console.print(f"[bold green]Task triggered with file_id:[/bold green] {file_id}")
        print("\n")
        subscribe_once(file_id)

if __name__ == "__main__":
    main()