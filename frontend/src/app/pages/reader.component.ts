import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { DomSanitizer, SafeHtml } from '@angular/platform-browser';
import { ApiService } from '../api.service';

@Component({
  selector: 'app-reader',
  standalone: true,
  imports: [CommonModule, RouterLink],
  templateUrl: './reader.component.html'
})
export class ReaderPageComponent {
  bookId = '';
  pageNo = 1;
  totalPages = 0;
  html: SafeHtml = '';
  error = '';
  pages: { pageNo: number; section?: string; label?: string }[] = [];
  chapters: { number: number; title: string; pageNoStart?: number }[] = [];
  navCollapsed = false;
  expanded: Record<number, boolean> = {};
  activeChapter?: number;
  section = '';
  label = '';
  audioUrl = '';
  videoUrl = '';

  constructor(private route: ActivatedRoute, private router: Router, private api: ApiService, private san: DomSanitizer) {
    this.route.paramMap.subscribe(params => {
      this.bookId = params.get('id') || '';
      const n = parseInt(params.get('n') || '1', 10);
      this.pageNo = isNaN(n) || n < 1 ? 1 : n;
      this.load();
    });
  }

  load() {
    this.error = '';
    this.api.listPages(this.bookId).subscribe({
      next: res => { this.pages = res?.data?.pages || []; this.totalPages = res?.data?.totalPages || this.totalPages; },
      error: () => { this.pages = []; }
    });
    this.api.listChapters(this.bookId).subscribe({
      next: res => { this.chapters = res?.data || []; this.updateActiveChapter(); },
      error: () => { this.chapters = []; this.updateActiveChapter(); }
    });
    this.api.getPage(this.bookId, this.pageNo).subscribe({
      next: res => {
        const d = res?.data || {};
        this.totalPages = d.totalPages || 0;
        this.html = this.san.bypassSecurityTrustHtml(d.html || '<p>No content</p>');
        this.section = d.section || '';
        this.label = d.label || '';
        this.audioUrl = d.audioUrl || '';
        this.videoUrl = d.videoUrl || '';
        this.updateActiveChapter();
      },
      error: err => {
        if (err?.status === 403 || err?.status === 401) this.error = 'You do not have access to read this book.';
        else if (err?.status === 404) this.error = 'Page not found.';
        else this.error = 'Unable to load page.';
      }
    });
  }

  goPrev() { if (this.pageNo > 1) this.router.navigate(['/book', this.bookId, 'page', this.pageNo - 1]); }
  goNext() { if (!this.totalPages || this.pageNo < this.totalPages) this.router.navigate(['/book', this.bookId, 'page', this.pageNo + 1]); }
  go(n: number | { pageNo: number }) {
    const p = typeof n === 'number' ? n : n.pageNo;
    this.router.navigate(['/book', this.bookId, 'page', p]);
  }
  btnClass(p: { pageNo: number }) { return this.pageNo === p.pageNo ? 'bg-blue-600 text-white' : 'bg-slate-800 hover:bg-slate-700'; }

  toggleNav() { this.navCollapsed = !this.navCollapsed; }
  toggleChapter(num: number) { this.expanded[num] = !this.expanded[num]; }
  pagesForChapter(num: number) {
    const ch = this.chapters.find(c => c.number === num);
    if (!ch || !ch.pageNoStart) return [] as { pageNo: number }[];
    const idx = this.chapters.findIndex(c => c.number === num);
    const next = idx >= 0 && idx + 1 < this.chapters.length ? this.chapters[idx + 1] : undefined;
    const end = next?.pageNoStart ? (next.pageNoStart - 1) : (this.totalPages || Number.MAX_SAFE_INTEGER);
    return this.pages.filter(p => p.pageNo >= (ch.pageNoStart as number) && p.pageNo <= end);
  }
  updateActiveChapter() {
    if (!this.chapters?.length) { this.activeChapter = undefined; return; }
    let active: number | undefined;
    for (let i = 0; i < this.chapters.length; i++) {
      const cur = this.chapters[i];
      const start = cur.pageNoStart || 1;
      const next = (i + 1 < this.chapters.length) ? (this.chapters[i+1].pageNoStart || Number.MAX_SAFE_INTEGER) : Number.MAX_SAFE_INTEGER;
      if (this.pageNo >= start && this.pageNo < next) { active = cur.number; break; }
    }
    this.activeChapter = active;
  }
  isActiveChapter(c: {number:number; pageNoStart?: number}) { return this.activeChapter === c.number; }
}
