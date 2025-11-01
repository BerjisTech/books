import { Component, OnInit } from '@angular/core';
import { RouterLink, RouterOutlet } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { environment } from '../environments/environment';
import { CommonModule } from '@angular/common';
import { AuthService } from './auth.service';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, RouterOutlet, RouterLink, FormsModule],
  templateUrl: './app.component.html'
})
export class AppComponent implements OnInit {
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
  constructor(private auth: AuthService) {}
  ngOnInit(): void {
    const w: any = (typeof window !== 'undefined') ? (window as any) : {};
    const stored = typeof localStorage !== 'undefined' ? localStorage.getItem('booksApiBase') : '';
    if (stored && !w.__BOOKS_API__) { w.__BOOKS_API__ = stored; }
    this.booksBase = w.__BOOKS_API__ || 'http://localhost:8088';
    this.newBooksBase = this.booksBase;
    // Check auth/session against Core API
    this.auth.isAuthed().subscribe(v => this.authed = !!v);
    this.auth.user().subscribe(u => this.userName = u?.name || null);
    this.auth.check();
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
}
