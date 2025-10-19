import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { DomSanitizer, SafeHtml } from '@angular/platform-browser';
import { ApiService } from '../api.service';

@Component({
  selector: 'app-reader',
  standalone: true,
  imports: [CommonModule, RouterLink],
  template: `
    <div class="flex items-center justify-between mb-2">
      <h2 class="m-0 text-xl font-semibold">Reader</h2>
      <a routerLink="/marketplace" class="text-blue-500 hover:underline">Back to Marketplace</a>
    </div>
    <div *ngIf="error" class="p-2 border border-rose-600 text-rose-100 rounded">{{ error }}</div>
    <div *ngIf="!error" class="grid grid-cols-12 gap-4">
      <aside class="col-span-12 md:col-span-3">
        <div class="flex items-center justify-between mb-2">
          <div class="font-semibold">Contents</div>
          <button (click)="toggleNav()" class="text-xs px-2 py-1 border rounded">{{ navCollapsed ? 'Expand' : 'Collapse' }}</button>
        </div>
        <ng-container *ngIf="!navCollapsed">
          <div class="font-semibold mb-1" *ngIf="chapters.length">Chapters</div>
          <ul class="mb-3">
            <li *ngFor="let c of chapters" class="mb-1">
              <button (click)="toggleChapter(c.number)" [ngClass]="{'font-semibold text-blue-400': isActiveChapter(c), 'text-slate-300': !isActiveChapter(c)}" class="text-left w-full hover:underline">
                {{ c.number }}. {{ c.title }}
              </button>
              <div *ngIf="expanded[c.number]" class="mt-1 flex flex-wrap gap-2 text-xs">
                <button *ngFor="let p of pagesForChapter(c.number)" (click)="go(p)" [class]="btnClass(p)" class="px-2 py-1 rounded border">
                  {{ p.pageNo }}
                </button>
              </div>
            </li>
          </ul>
          <div class="font-semibold mb-1">All Pages</div>
          <div class="flex flex-wrap gap-2 text-sm">
            <button *ngFor="let p of pages" (click)="go(p)" [class]="btnClass(p)" class="px-2 py-1 rounded border">
              {{ p.pageNo }}
            </button>
          </div>
        </ng-container>
      </aside>
      <section class="col-span-12 md:col-span-9">
        <div class="mb-2">
          <button (click)="goPrev()" [disabled]="pageNo<=1" class="px-2 py-1 border rounded disabled:opacity-50">Prev</button>
          <span class="mx-2">Page {{ pageNo }} of {{ totalPages || '?' }}</span>
          <button (click)="goNext()" [disabled]="totalPages>0 && pageNo>=totalPages" class="px-2 py-1 border rounded disabled:opacity-50">Next</button>
        </div>
        <div [innerHTML]="html" class="min-h-[60vh] p-4 border rounded bg-slate-900"></div>
      </section>
    </div>
  `
})
export class ReaderPageComponent {
  bookId = '';
  pageNo = 1;
  totalPages = 0;
  html: SafeHtml = '';
  error = '';
  pages: { pageNo: number }[] = [];
  chapters: { number: number; title: string; pageNoStart?: number }[] = [];
  navCollapsed = false;
  expanded: Record<number, boolean> = {};
  activeChapter?: number;

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
  go(n: number | {pageNo: number}) {
    const p = typeof n === 'number' ? n : n.pageNo;
    this.router.navigate(['/book', this.bookId, 'page', p]);
  }
  btnClass(p: {pageNo: number}) { return this.pageNo === p.pageNo ? 'bg-blue-600 text-white' : 'bg-slate-800 hover:bg-slate-700'; }

  toggleNav() { this.navCollapsed = !this.navCollapsed; }
  toggleChapter(num: number) { this.expanded[num] = !this.expanded[num]; }
  pagesForChapter(num: number) {
    const ch = this.chapters.find(c => c.number === num);
    if (!ch || !ch.pageNoStart) return [] as {pageNo:number}[];
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
