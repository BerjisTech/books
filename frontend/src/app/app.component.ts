import { Component, OnInit } from '@angular/core';
import { RouterLink, RouterOutlet } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { environment } from '../environments/environment';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, RouterLink, FormsModule],
  templateUrl: './app.component.html'
})
export class AppComponent implements OnInit {
  isDev = !environment.production;
  coreBase = environment.apiBase;
  booksBase = '';
  newBooksBase = '';
  ngOnInit(): void {
    const w: any = (typeof window !== 'undefined') ? (window as any) : {};
    const stored = typeof localStorage !== 'undefined' ? localStorage.getItem('booksApiBase') : '';
    if (stored && !w.__BOOKS_API__) { w.__BOOKS_API__ = stored; }
    this.booksBase = w.__BOOKS_API__ || 'http://localhost:8088';
    this.newBooksBase = this.booksBase;
  }
  saveBooksBase() {
    const w: any = (typeof window !== 'undefined') ? (window as any) : {};
    if (this.newBooksBase && this.newBooksBase.trim().length > 0) {
      w.__BOOKS_API__ = this.newBooksBase.trim();
      try { localStorage.setItem('booksApiBase', w.__BOOKS_API__); } catch {}
      this.booksBase = w.__BOOKS_API__;
      alert('Books API base set to: ' + this.booksBase);
    } else {
      delete w.__BOOKS_API__;
      try { localStorage.removeItem('booksApiBase'); } catch {}
      this.booksBase = 'http://localhost:8088';
      alert('Books API base cleared; using default ' + this.booksBase);
    }
  }
}
