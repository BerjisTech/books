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
  template: `
    <h2>Dashboard</h2>

    <div *ngIf="!authed" style="padding:8px; border:1px solid #e5e7eb22; border-radius:6px; background:#0b1020;">
      <div style="font-weight:600;">You’re not signed in</div>
      <p style="opacity:.85;">View the marketplace publicly. To manage your profile, purchases, or publish books, sign in at <a href="http://berjis.test/account">berjis.test</a>.</p>
    </div>

    <div *ngIf="authed" style="margin:8px 0;">
      <div style="margin-bottom:8px;">Roles: <span *ngIf="roles.length === 0">reader (default)</span><span *ngFor="let r of roles" style="display:inline-block; margin-right:6px; padding:2px 8px; border:1px solid #e5e7eb22; border-radius:999px; font-size:12px;">{{ r }}</span></div>

      <section style="margin:16px 0;">
        <h3 style="margin:0 0 6px 0;">Reader</h3>
        <p style="opacity:.85;">Your library and recent purchases.</p>
        <div *ngIf="libLoading">Loading library…</div>
        <div *ngIf="!libLoading && library.length === 0" style="opacity:.8">No purchases yet.</div>
        <ul>
          <li *ngFor="let item of library" style="margin:6px 0;">
            <strong>{{ item.title || item.bookId }}</strong>
            <span style="opacity:.8"> — {{ item.kind }} • {{ item.status }}</span>
            <button *ngIf="item.kind==='ebook' && item.status==='completed'" (click)="read(item.bookId)" style="margin-left:8px;">Read</button>
          </li>
        </ul>
      </section>

      <section *ngIf="hasRole('author') || hasRole('creator')" style="margin:16px 0;">
        <h3 style="margin:0 0 6px 0;">Author</h3>
        <p style="opacity:.85;">Create and manage your books.</p>
        <div style="margin:8px 0;">
          <input [(ngModel)]="newTitle" placeholder="Book title" />
          <input [(ngModel)]="newDesc" placeholder="Description (optional)" style="width:320px;" />
          <button (click)="createBook()">Create</button>
          <span *ngIf="createMsg" style="margin-left:8px; font-size:12px; opacity:.85;">{{ createMsg }}</span>
        </div>
        <div style="margin-top:8px;">
          <div style="font-weight:600; margin:8px 0 4px 0;">My Books</div>
          <div *ngIf="myLoading">Loading…</div>
          <div *ngIf="!myLoading && myBooks.length===0" style="opacity:.8">No books yet.</div>
          <ul>
            <li *ngFor="let b of myBooks" style="margin:6px 0;">
              <div>
                <strong>{{ b.title }}</strong>
                <span *ngIf="b.subtitle"> — {{ b.subtitle }}</span>
                <span style="opacity:.8"> • {{ b.visibility }}</span>
                <button style="margin-left:8px;" (click)="toggleEditor(b.id)">{{ editBookId===b.id ? 'Close' : 'Add Page' }}</button>
                <button style="margin-left:8px;" (click)="toggleChapters(b.id)">{{ chapterBookId===b.id ? 'Hide Chapters' : 'Chapters' }}</button>
                <a [routerLink]="['/book', b.id]" style="margin-left:8px;">Open Reader</a>
              </div>
              <div *ngIf="editBookId===b.id" style="margin-top:6px; padding:8px; border:1px solid #e5e7eb22; border-radius:6px;">
                <div>
                  Page #: <input type="number" [(ngModel)]="newPageNo" min="1" style="width:80px;" />
                </div>
                <div style="margin-top:6px;">
                  <div style="display:flex; gap:8px; align-items:center; margin-bottom:6px;">
                    <label style="font-size:12px; opacity:.8;">
                      <input type="checkbox" [(ngModel)]="useMarkdown" /> Markdown
                    </label>
                    <label *ngIf="useMarkdown" style="font-size:12px; opacity:.8;">
                      <input type="checkbox" [(ngModel)]="previewMd" /> Preview
                    </label>
                    <button (click)="insertImage(mdArea)" style="font-size:12px;">Insert Image</button>
                    <button (click)="insertCode(mdArea)" style="font-size:12px;">Insert Code</button>
                  </div>
                  <textarea #mdArea [(ngModel)]="newPageHtml" rows="6" style="width:100%;" placeholder="Page HTML or Markdown"></textarea>
                  <div *ngIf="useMarkdown && previewMd" style="margin-top:8px; padding:8px; border:1px dashed #e5e7eb33; border-radius:6px; background:#0f172a;">
                    <div style="font-size:12px; opacity:.8; margin-bottom:4px;">Preview</div>
                    <div [innerHTML]="mdToHtml(newPageHtml)"></div>
                  </div>
                </div>
                <div style="margin-top:6px;">
                  <button (click)="savePage(b.id)">Save Page</button>
                </div>
              </div>
              <div *ngIf="chapterBookId===b.id" style="margin-top:6px; padding:8px; border:1px solid #e5e7eb22; border-radius:6px;">
                <div style="font-weight:600; margin-bottom:6px;">Chapters</div>
                <div *ngIf="chapLoading">Loading chapters…</div>
                <ul>
                  <li *ngFor="let c of chapters" style="font-size:14px; margin:2px 0;">{{ c.number }}. {{ c.title }} <span *ngIf="c.pageNoStart" style="opacity:.7">(p{{ c.pageNoStart }})</span></li>
                </ul>
                <div style="margin-top:8px; font-size:13px;">Add/Update chapter</div>
                <div>Number: <input type="number" [(ngModel)]="chNumber" min="1" style="width:80px;" /> Title: <input [(ngModel)]="chTitle" style="width:240px;" /> Start Page: <input type="number" [(ngModel)]="chStart" min="1" style="width:80px;" /></div>
                <div style="margin-top:6px;"><button (click)="saveChapter(b.id)">Save Chapter</button></div>
              </div>
            </li>
          </ul>
        </div>
      </section>

      <section *ngIf="hasRole('publisher')" style="margin:16px 0;">
        <h3 style="margin:0 0 6px 0;">Publisher</h3>
        <p style="opacity:.85;">Organize catalog, pricing, and share permissions. (UI placeholder)</p>
      </section>

      <section *ngIf="hasRole('organizer') || hasRole('club_admin')" style="margin:16px 0;">
        <h3 style="margin:0 0 6px 0;">Community</h3>
        <p style="opacity:.85;">Manage clubs and meetups. (UI placeholder)</p>
      </section>
    </div>
  `
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
  // Chapters state
  chapterBookId: string | null = null;
  chapters: any[] = [];
  chapLoading = false;
  chNumber: number = 1;
  chTitle: string = '';
  chStart?: number;

  constructor(private auth: AuthService, private api: ApiService, private router: Router) {}

  ngOnInit(): void {
    this.auth.isAuthed().subscribe(a => { this.authed = !!a; if (this.authed) { this.loadLibrary(); this.loadMyBooks(); } });
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
    this.createMsg = 'Creating…';
    this.api.createBook(t, this.newDesc || undefined).subscribe({
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
