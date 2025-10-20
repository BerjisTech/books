import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { AuthService } from '../auth.service';
import { ApiService } from '../api.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './dashboard.component.html'
})
export class DashboardPageComponent implements OnInit {
  authed = false;
  roles: string[] = [];
  library: any[] = [];
  libLoading = false;
  createMsg = '';
  newTitle = '';
  newDesc = '';
  myBooks: any[] = [];
  myLoading = false;
  editBookId: string | null = null;
  newPageNo: number = 1;
  newPageHtml: string = '';
  useMarkdown = true;
  previewMd = false;
  // Profiles
  myAuthors: any[] = [];
  myPublishers: any[] = [];
  profileKind: 'author'|'publisher' = 'author';
  selectedAuthorId: string | null = null;
  selectedPublisherId: string | null = null;
  // Chapters state
  chapterBookId: string | null = null;
  chapters: any[] = [];
  chapLoading = false;
  chNumber: number = 1;
  chTitle: string = '';
  chStart?: number;

  constructor(private auth: AuthService, private api: ApiService, private router: Router) {}

  ngOnInit(): void {
    this.auth.isAuthed().subscribe(a => { this.authed = !!a; if (this.authed) { this.loadLibrary(); this.loadMyBooks(); this.loadProfiles(); } });
    this.auth.roles().subscribe(r => this.roles = r || []);
    // ensure session checked on entry
    this.auth.check();
  }

  hasRole(r: string) { return this.roles.includes(r); }

  loadLibrary() {
    this.libLoading = true;
    this.api.listMyPurchases().subscribe({
      next: res => { this.library = res.data || []; this.libLoading = false; },
      error: () => { this.library = []; this.libLoading = false; }
    });
  }

  createBook() {
    const t = this.newTitle.trim();
    if (!t) { this.createMsg = 'Title is required'; return; }
    const payload: any = { title: t };
    if (this.newDesc) payload.description = this.newDesc;
    if (this.profileKind === 'author' && this.selectedAuthorId) payload.authorId = this.selectedAuthorId;
    if (this.profileKind === 'publisher' && this.selectedPublisherId) payload.publisherId = this.selectedPublisherId;
    this.createMsg = 'Creating…';
    this.api.createBook(payload).subscribe({
      next: res => { this.createMsg = 'Created'; this.newTitle = ''; this.newDesc = ''; this.loadMyBooks(); },
      error: () => { this.createMsg = 'Failed'; }
    });
  }

  loadMyBooks() {
    this.myLoading = true;
    this.api.listMyBooks().subscribe({
      next: res => { this.myBooks = res.data || []; this.myLoading = false; },
      error: () => { this.myBooks = []; this.myLoading = false; }
    });
  }

  loadProfiles() {
    this.api.listMyAuthors().subscribe({ next: res => { this.myAuthors = res.data || []; }, error: () => this.myAuthors = [] });
    this.api.listMyPublishers().subscribe({ next: res => { this.myPublishers = res.data || []; }, error: () => this.myPublishers = [] });
  }
  applyAuthor() {
    const name = prompt('Display name for author profile');
    if (!name) return;
    this.api.applyAuthor(name).subscribe({ next: () => this.loadProfiles(), error: () => alert('Failed to apply as author') });
  }
  applyPublisher() {
    const name = prompt('Publisher name');
    if (!name) return;
    this.api.applyPublisher(name).subscribe({ next: () => this.loadProfiles(), error: () => alert('Failed to apply as publisher') });
  }

  read(bookId: string) { this.router.navigate(['/book', bookId]); }

  toggleEditor(bookId: string) {
    if (this.editBookId === bookId) { this.editBookId = null; return; }
    this.editBookId = bookId;
    this.newPageNo = 1; this.newPageHtml = '';
  }
  toggleChapters(bookId: string) {
    if (this.chapterBookId === bookId) { this.chapterBookId = null; return; }
    this.chapterBookId = bookId; this.loadChapters(bookId);
  }
  loadChapters(bookId: string) {
    this.chapLoading = true;
    this.api.listChapters(bookId).subscribe({
      next: res => { this.chapters = res.data || []; this.chapLoading = false; },
      error: () => { this.chapters = []; this.chapLoading = false; }
    });
  }
  saveChapter(bookId: string) {
    if (!this.chNumber || !this.chTitle.trim()) { alert('Chapter number and title required'); return; }
    this.api.upsertChapter(bookId, this.chNumber, this.chTitle.trim(), this.chStart || undefined).subscribe({
      next: () => { this.loadChapters(bookId); this.chTitle=''; },
      error: () => alert('Failed to save chapter')
    });
  }
  savePage(bookId: string) {
    if (!this.newPageNo || this.newPageNo < 1) { alert('Page number must be >= 1'); return; }
    const html = this.useMarkdown ? this.mdToHtml(this.newPageHtml || '') : (this.newPageHtml || '');
    this.api.upsertPage(bookId, this.newPageNo, html).subscribe({
      next: () => { alert('Saved'); this.newPageHtml = ''; },
      error: () => alert('Failed to save page')
    });
  }

  // Minimal Markdown -> HTML converter (headings, bold, italic, lists, links, paragraphs)
  mdToHtml(md: string): string {
    let out = md.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
    // headings
    out = out.replace(/^###\s+(.+)$/gm, '<h3>$1</h3>');
    out = out.replace(/^##\s+(.+)$/gm, '<h2>$1</h2>');
    out = out.replace(/^#\s+(.+)$/gm, '<h1>$1</h1>');
    // bold/italic
    out = out.replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>');
    out = out.replace(/\*(.+?)\*/g, '<em>$1</em>');
    // links [text](url)
    out = out.replace(/\[(.+?)\]\((https?:\/\/[^\s)]+)\)/g, '<a href="$2" target="_blank" rel="noopener">$1<\/a>');
    // unordered lists
    out = out.replace(/(?:^|\n)(?:-\s+.+\n?)+/g, (block) => {
      const items = block.trim().split(/\n/).map(l => l.replace(/^[-*]\s+/, '')).filter(Boolean);
      return '<ul>' + items.map(i => `<li>${i}</li>`).join('') + '</ul>';
    });
    // paragraphs
    out = out.split(/\n{2,}/).map(p => {
      if (/^<h\d|^<ul/.test(p.trim())) return p; // already block
      const lines = p.split(/\n/).filter(Boolean).join('<br/>');
      return `<p>${lines}</p>`;
    }).join('\n');
    return out;
  }

  insertAtCursor(textarea: HTMLTextAreaElement, snippet: string) {
    if (!textarea) { this.newPageHtml += snippet; return; }
    const start = textarea.selectionStart || 0;
    const end = textarea.selectionEnd || 0;
    const before = this.newPageHtml.slice(0, start);
    const after = this.newPageHtml.slice(end);
    this.newPageHtml = before + snippet + after;
    setTimeout(() => {
      const pos = start + snippet.length;
      try { textarea.setSelectionRange(pos, pos); textarea.focus(); } catch {}
    });
  }
  insertImage(el: HTMLTextAreaElement) {
    const url = prompt('Image URL');
    const alt = prompt('Alt text') || 'image';
    if (url) this.insertAtCursor(el, `![${alt}](${url})`);
  }
  insertCode(el: HTMLTextAreaElement) {
    const lang = prompt('Language (optional)') || '';
    const code = 'your code here';
    const fence = '```' + lang + '\n' + code + '\n```\n';
    this.insertAtCursor(el, fence);
  }
}
