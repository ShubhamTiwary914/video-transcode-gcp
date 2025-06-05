from rich import print
from rich.console import Console
from rich.progress import Progress, SpinnerColumn, TextColumn, BarColumn, TimeElapsedColumn
from subprocess import run, PIPE
from tkinter import Tk, filedialog
from task_trigger import trigger_task
import os
import sys
import shlex

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

    console.print("[bold green]Signed URL:[/bold green]", download_url)
    print("\n")
    console.print("[cyan]Triggering processing task...[/cyan]")
    file_id = trigger_task(download_url)
    if file_id:
        console.print(f"[bold green]Task triggered with file_id:[/bold green] {file_id}")

if __name__ == "__main__":
    main()
