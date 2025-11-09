import { Component, OnDestroy, OnInit } from '@angular/core';
import { RouterLink, RouterOutlet } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { environment } from '../environments/environment';
import { CommonModule } from '@angular/common';
import { CoreAuthService, CoreAuthSession } from '@berjis/angular-auth';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, RouterOutlet, RouterLink, FormsModule],
  templateUrl: './app.component.html'
})
export class AppComponent implements OnInit, OnDestroy {
  isDark = false;
  isDev = !environment.production;
  coreBase = environment.apiBase;
  booksBase = '';
  newBooksBase = '';
  authed = false;
  accountUrl = 'http://berjis.tech/account';
  userName: string | null = null;
  devNotice: { text: string; tone: 'info' | 'success' | 'error' } | null = null;
  collectionOpen = false;
  marketplaceOpen = false;
  communityOpen = false;
  private unsubscribeAuth?: () => void;

  constructor(private auth: CoreAuthService) {}
  ngOnInit(): void {
    const persisted = (localStorage.getItem('theme') || '').toLowerCase();
    const preferDark = persisted === 'dark';
    this.setTheme(preferDark ? 'dark' : 'light');
    const w: any = (typeof window !== 'undefined') ? (window as any) : {};
    const stored = typeof localStorage !== 'undefined' ? localStorage.getItem('booksApiBase') : '';
    if (stored && !w.__BOOKS_API__) { w.__BOOKS_API__ = stored; }
    this.booksBase = w.__BOOKS_API__ || 'http://localhost:8088';
    this.newBooksBase = this.booksBase;
    this.unsubscribeAuth = this.auth.onSessionChange((session: CoreAuthSession) => this.applySession(session));
    void this.auth.ensureAuth({ maxAgeMs: 1500 }).catch(() => {});
  }

  ngOnDestroy(): void {
    this.unsubscribeAuth?.();
  }
  toggleTheme() { this.setTheme(this.isDark ? 'light' : 'dark'); }
  private setTheme(mode: 'light' | 'dark') {
    this.isDark = mode === 'dark';
    document.documentElement.classList.toggle('dark', mode === 'dark');
    try { localStorage.setItem('theme', mode); } catch {}
  }
  toggleCollection() { this.collectionOpen = !this.collectionOpen; }
  toggleMarketplace() { this.marketplaceOpen = !this.marketplaceOpen; }
  toggleCommunity() { this.communityOpen = !this.communityOpen; }
  private setDevNotice(text: string, tone: 'info' | 'success' | 'error' = 'info') {
    this.devNotice = { text, tone };
    setTimeout(() => { this.devNotice = null; }, 5000);
  }
  saveBooksBase() {
    const w: any = (typeof window !== 'undefined') ? (window as any) : {};
    if (this.newBooksBase && this.newBooksBase.trim().length > 0) {
      w.__BOOKS_API__ = this.newBooksBase.trim();
      try { localStorage.setItem('booksApiBase', w.__BOOKS_API__); } catch {}
      this.booksBase = w.__BOOKS_API__;
      this.setDevNotice('Books API base updated.', 'success');
    } else {
      delete w.__BOOKS_API__;
      try { localStorage.removeItem('booksApiBase'); } catch {}
      this.booksBase = 'http://localhost:8088';
      this.setDevNotice('Books API base cleared; using default.', 'info');
    }
  }

  private applySession(session: CoreAuthSession) {
    this.authed = !!session?.valid;
    this.userName = this.resolveDisplayName(session);
  }

  private resolveDisplayName(session: CoreAuthSession): string | null {
    if (!session) {
      return null;
    }
    const profile = session.profile ?? {};
    const candidates: Array<unknown> = [
      (profile as Record<string, unknown>)['displayName'],
      (profile as Record<string, unknown>)['name'],
      (profile as Record<string, unknown>)['fullName'],
      (profile as Record<string, unknown>)['username'],
      session.email
    ];
    for (const value of candidates) {
      if (typeof value === 'string' && value.trim().length > 0) {
        return value.trim();
      }
    }
    return null;
  }
}
