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
  newPageSection: 'frontmatter' | 'content' | 'backmatter' = 'content';
  newPageLabel: string = '';
  newPageAudio: string = '';
  newPageVideo: string = '';
  assetsByBook: Record<string, any[]> = {};
  assetForm: Record<string, { kind: string; source: string; url: string; checksum?: string; fileSizeBytes?: number | string }> = {};
  assetPanelBookId: string | null = null;
  editorArea: HTMLTextAreaElement | null = null;
  showImageHelper = false;
  showCodeHelper = false;
  imageHelper = { url: '', alt: '' };
  codeHelper = { lang: '', snippet: 'your code here' };
  collaborations: any[] = [];
  collabLoading = false;
  earnings: { entries: any[]; totals: Record<string, number> } | null = null;
  earningsLoading = false;
  notices: Record<string, { text: string; tone: 'info' | 'success' | 'error' } | null> = {
    creator: null,
    pages: null,
    chapters: null,
    assets: null,
    collaborations: null,
    earnings: null,
    general: null
  };

  constructor(private auth: AuthService, private api: ApiService, private router: Router) {}

  private setNotice(key: keyof typeof this.notices, text: string, tone: 'info' | 'success' | 'error' = 'info') {
    this.notices[key] = { text, tone };
  }

  clearNotice(key: keyof typeof this.notices) {
    this.notices[key] = null;
  }

  ngOnInit(): void {
    this.auth.isAuthed().subscribe(a => {
      this.authed = !!a;
      if (this.authed) {
        this.loadLibrary();
        this.loadMyBooks();
        this.loadProfiles();
        this.loadCollaborations();
        this.loadEarnings();
      }
    });
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
    if (!t) {
      this.setNotice('creator', 'Add a working title to create a draft.', 'error');
      return;
    }
    const payload: any = { title: t };
    if (this.newDesc) payload.description = this.newDesc;
    if (this.profileKind === 'author' && this.selectedAuthorId) payload.authorId = this.selectedAuthorId;
    if (this.profileKind === 'publisher' && this.selectedPublisherId) payload.publisherId = this.selectedPublisherId;
    this.setNotice('creator', 'Creating draft…', 'info');
    this.api.createBook(payload).subscribe({
      next: () => {
        this.setNotice('creator', 'Draft created. Add pages whenever you are ready.', 'success');
        this.newTitle = '';
        this.newDesc = '';
        this.loadMyBooks();
      },
      error: () => {
        this.setNotice('creator', 'We could not create the draft right now. Please try again.', 'error');
      }
    });
  }

  loadMyBooks() {
    this.myLoading = true;
    this.assetsByBook = {};
    this.assetForm = {};
    this.api.listMyBooks().subscribe({
      next: res => { this.myBooks = res.data || []; this.myLoading = false; },
      error: () => { this.myBooks = []; this.myLoading = false; }
    });
  }

  loadProfiles() {
    this.api.listMyAuthors().subscribe({ next: res => { this.myAuthors = res.data || []; }, error: () => this.myAuthors = [] });
    this.api.listMyPublishers().subscribe({ next: res => { this.myPublishers = res.data || []; }, error: () => this.myPublishers = [] });
  }


  loadCollaborations() {
    this.collabLoading = true;
    this.api.listCollaborations().subscribe({
      next: res => {
        this.collaborations = res.data || [];
        this.collabLoading = false;
        this.clearNotice('collaborations');
      },
      error: () => {
        this.collaborations = [];
        this.collabLoading = false;
        this.setNotice('collaborations', 'Unable to load collaboration requests right now.', 'error');
      }
    });
  }
  respondCollaboration(collab: any, decision: 'accept' | 'decline') {
    this.api.respondCollaboration(collab.id, decision).subscribe({
      next: () => {
        this.loadCollaborations();
        this.loadMyBooks();
        this.setNotice('collaborations', decision === 'accept' ? 'Collaboration accepted.' : 'Collaboration declined.', 'success');
      },
      error: () => this.setNotice('collaborations', 'We could not update that collaboration right now.', 'error')
    });
  }

  loadEarnings() {
    this.earningsLoading = true;
    this.api.listEarnings().subscribe({
      next: res => {
        this.earnings = res.data || { entries: [], totals: {} };
        this.earningsLoading = false;
        this.clearNotice('earnings');
      },
      error: () => {
        this.earnings = { entries: [], totals: {} };
        this.earningsLoading = false;
        this.setNotice('earnings', 'Unable to load earnings right now.', 'error');
      }
    });
  }

  read(bookId: string) { this.router.navigate(['/book', bookId]); }

  toggleEditor(bookId: string) {
    if (this.editBookId === bookId) { this.editBookId = null; return; }
    this.editBookId = bookId;
    this.newPageNo = 1;
    this.newPageHtml = '';
    this.newPageSection = 'content';
    this.newPageLabel = '';
    this.newPageAudio = '';
    this.newPageVideo = '';
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
    if (!this.chNumber || !this.chTitle.trim()) { this.setNotice('chapters', 'Provide both a chapter number and a title.', 'error'); return; }
    this.api.upsertChapter(bookId, {
      number: this.chNumber,
      title: this.chTitle.trim(),
      pageNoStart: this.chStart || undefined
    }).subscribe({
      next: () => { this.loadChapters(bookId); this.chTitle=''; },
      error: () => this.setNotice('chapters', 'We could not save that chapter right now.', 'error')
    });
  }
  toggleAssets(bookId: string) {
    if (this.assetPanelBookId === bookId) { this.assetPanelBookId = null; return; }
    this.assetPanelBookId = bookId;
    if (!this.assetForm[bookId]) this.assetForm[bookId] = { kind: 'ebook', source: 'external', url: '' };
    this.loadAssets(bookId);
  }
  loadAssets(bookId: string) {
    this.api.listBookAssets(bookId).subscribe({
      next: res => { this.assetsByBook[bookId] = res.data || []; },
      error: () => { this.assetsByBook[bookId] = []; }
    });
  }
  addAsset(bookId: string) {
    const form = this.assetForm[bookId] || { kind: 'ebook', source: 'external', url: '' };
    if (!form.url || !form.url.trim()) { this.setNotice('assets', 'Add an asset URL before saving.', 'error'); return; }
    const payload: any = {
      kind: form.kind || 'ebook',
      source: form.source || 'external',
      url: form.url.trim()
    };
    if (form.checksum) payload.checksum = form.checksum;
    if (form.fileSizeBytes !== undefined && form.fileSizeBytes !== null && form.fileSizeBytes !== '') {
      const size = Number(form.fileSizeBytes);
      if (!Number.isNaN(size)) payload.fileSizeBytes = size;
    }
    this.api.createBookAsset(bookId, payload).subscribe({
      next: () => {
        this.assetForm[bookId] = { kind: form.kind, source: form.source, url: '', checksum: '', fileSizeBytes: undefined };
        this.loadAssets(bookId);
      },
      error: () => this.setNotice('assets', 'Unable to attach that asset right now.', 'error')
    });
  }
  removeAsset(bookId: string, assetId: string) {
    this.api.deleteBookAsset(bookId, assetId).subscribe({
      next: () => {
        this.setNotice('assets', 'Asset removed from the draft.', 'success');
        this.loadAssets(bookId);
      },
      error: () => this.setNotice('assets', 'Unable to remove that asset right now.', 'error')
    });
  }
  publishBook(book: any) {
    this.api.publishBook(book.id, book?.releaseDate ? book.releaseDate.substring(0, 10) : undefined).subscribe({
      next: () => {
        this.setNotice('creator', 'Book published.', 'success');
        this.loadMyBooks();
      },
      error: () => this.setNotice('creator', 'We could not publish that book right now.', 'error')
    });
  }
  unpublishBook(book: any) {
    this.api.unpublishBook(book.id).subscribe({
      next: () => {
        this.setNotice('creator', 'Book returned to draft.', 'info');
        this.loadMyBooks();
      },
      error: () => this.setNotice('creator', 'We could not unpublish that book right now.', 'error')
    });
  }
  formatCurrency(amount?: number, currency?: string) {
    if (!amount) return 'Free';
    const cur = currency || 'USD';
    return `${amount.toFixed(2)} ${cur}`;
  }
  savePage(bookId: string) {
    if (!this.newPageNo || this.newPageNo < 1) {
      this.setNotice('pages', 'Page number must be at least 1.', 'error');
      return;
    }
    const html = this.useMarkdown ? this.mdToHtml(this.newPageHtml || '') : (this.newPageHtml || '');
    this.api.upsertPage(bookId, {
      pageNo: this.newPageNo,
      html,
      section: this.newPageSection,
      label: this.newPageLabel || undefined,
      audioUrl: this.newPageAudio || undefined,
      videoUrl: this.newPageVideo || undefined
    }).subscribe({
      next: () => {
        this.setNotice('pages', 'Page saved.', 'success');
        this.newPageHtml = '';
        this.newPageLabel = '';
        this.newPageAudio = '';
        this.newPageVideo = '';
      },
      error: () => this.setNotice('pages', 'We could not save that page right now.', 'error')
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

  insertAtCursor(textarea: HTMLTextAreaElement | null, snippet: string) {
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
    this.editorArea = el;
    this.showImageHelper = !this.showImageHelper;
    if (this.showImageHelper) {
      this.showCodeHelper = false;
    }
  }
  insertCode(el: HTMLTextAreaElement) {
    this.editorArea = el;
    this.showCodeHelper = !this.showCodeHelper;
    if (this.showCodeHelper) {
      this.showImageHelper = false;
    }
  }
  applyImage() {
    const url = this.imageHelper.url.trim();
    if (!url) {
      this.setNotice('pages', 'Add an image URL before inserting.', 'error');
      return;
    }
    const alt = this.imageHelper.alt.trim() || 'image';
    this.insertAtCursor(this.editorArea, `![${alt}](${url})`);
    this.showImageHelper = false;
    this.imageHelper = { url: '', alt: '' };
    this.setNotice('pages', 'Image markup added to the page.', 'success');
    try { this.editorArea?.focus(); } catch {}
  }
  cancelImage() {
    this.showImageHelper = false;
    this.imageHelper = { url: '', alt: '' };
  }
  applyCode() {
    const snippet = (this.codeHelper.snippet || 'your code here').trim() || 'your code here';
    const lang = this.codeHelper.lang.trim();
    const fence = `\`\`\`${lang ? lang + '\n' : '\n'}${snippet}\n\`\`\`\n`;
    this.insertAtCursor(this.editorArea, fence);
    this.showCodeHelper = false;
    this.codeHelper = { lang, snippet: 'your code here' };
    this.setNotice('pages', 'Code block added to the page.', 'success');
    try { this.editorArea?.focus(); } catch {}
  }
  cancelCode() {
    this.showCodeHelper = false;
    this.codeHelper = { lang: '', snippet: 'your code here' };
  }
}
