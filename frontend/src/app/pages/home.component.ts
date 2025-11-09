import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { Book } from '../interfaces/book';
import { BooksComponent } from "../components/books/books.component";
import { BooksService } from '../services/books.service';
import { CommonModule } from '@angular/common';
import { BookCardComponent } from "../components/book-card/book-card.component";
import { CoreAuthService, CoreAuthSession } from '@berjis/angular-auth';

@Component({
  selector: 'app-home',
  standalone: true,
  templateUrl: './home.component.html',
  imports: [CommonModule, BooksComponent]
})
export class HomePageComponent implements OnInit, OnDestroy {

  authed = false;
  accountUrl = 'http://berjis.tech/account';
  userName: string | null = null;
  pageTitle: string = 'Book Recomendations';
  books: Book[] = []

  private unsubscribeAuth?: () => void;

  constructor(private auth: CoreAuthService) { }

  ngOnInit(): void {
    this.unsubscribeAuth = this.auth.onSessionChange((session: CoreAuthSession) => {
      this.authed = !!session?.valid;
      this.userName = this.resolveDisplayName(session);
    });
    void this.auth.ensureAuth({ maxAgeMs: 1500 }).catch(() => {});
  }

  ngOnDestroy(): void {
    this.unsubscribeAuth?.();
  }

  private resolveDisplayName(session: CoreAuthSession | null | undefined): string | null {
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
