export function Footer() {
  return (
    <footer className="border-t border-parchment-dark px-4 py-6">
      <div className="mx-auto flex max-w-6xl flex-col items-center justify-between gap-2 text-sm text-feather sm:flex-row">
        <span>&copy; 2024 Stanza Bonanza</span>
        <div className="flex gap-4">
          <a href="#" className="transition-colors hover:text-ink">Terms</a>
          <a href="#" className="transition-colors hover:text-ink">Privacy</a>
        </div>
      </div>
    </footer>
  );
}
