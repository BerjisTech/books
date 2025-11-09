import 'zone.js';
// Initialize Books API base from localStorage before app boot (dev convenience)
try {
  if (typeof window !== 'undefined') {
    const w: any = window as any;
    const stored = localStorage.getItem('booksApiBase');
    if (stored) { w.__BOOKS_API__ = stored; }
  }
} catch {}
import { bootstrapApplication } from '@angular/platform-browser';
import { provideHttpClient } from '@angular/common/http';
import { provideRouter, Routes } from '@angular/router';
import { AppComponent } from './app/app.component';
import { HomePageComponent } from './app/pages/home.component';
import { MarketplacePageComponent } from './app/pages/marketplace.component';
import { ClubsPageComponent } from './app/pages/clubs.component';
import { MeetupsPageComponent } from './app/pages/meetups.component';
import { OrdersPageComponent } from './app/pages/orders.component';
import { DashboardPageComponent } from './app/pages/dashboard.component';
import { ReaderPageComponent } from './app/pages/reader.component';
import { BookComponent } from './app/components/book/book.component';
import { AuthorsPageComponent } from './app/pages/authors.component';
import { AuthorEnrollmentComponent } from './app/pages/author-enrollment.component';
import { PublisherEnrollmentComponent } from './app/pages/publisher-enrollment.component';
import { CORE_AUTH_API_BASE, createAuthGuard } from '@berjis/angular-auth';
import { environment } from './environments/environment';

const authGuard = createAuthGuard({
  ensureOptions: { maxAgeMs: 1500 }
});

const routes: Routes = [
  { path: '', component: HomePageComponent },
  { path: 'marketplace', component: MarketplacePageComponent },
  { path: 'clubs', component: ClubsPageComponent },
  { path: 'meetups', component: MeetupsPageComponent },
  { path: 'orders', component: OrdersPageComponent, canActivate: [authGuard] },
  { path: 'book/:id', component: BookComponent },
  { path: 'book/:id/page/:n', component: ReaderPageComponent },
  { path: 'authors', component: AuthorsPageComponent },
  { path: 'enroll/author', component: AuthorEnrollmentComponent, canActivate: [authGuard] },
  { path: 'enroll/publisher', component: PublisherEnrollmentComponent, canActivate: [authGuard] },
  { path: 'dashboard', component: DashboardPageComponent, canActivate: [authGuard] },
  { path: '**', redirectTo: '' }
];

bootstrapApplication(AppComponent, {
  providers: [
    provideHttpClient(),
    provideRouter(routes),
    { provide: CORE_AUTH_API_BASE, useValue: environment.apiBase }
  ]
}).catch(err => console.error(err));
