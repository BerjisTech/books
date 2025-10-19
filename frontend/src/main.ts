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

const routes: Routes = [
  { path: '', component: HomePageComponent },
  { path: 'marketplace', component: MarketplacePageComponent },
  { path: 'clubs', component: ClubsPageComponent },
  { path: 'meetups', component: MeetupsPageComponent },
  { path: 'orders', component: OrdersPageComponent },
  { path: 'book/:id', redirectTo: 'book/:id/page/1' },
  { path: 'book/:id/page/:n', component: ReaderPageComponent },
  { path: 'dashboard', component: DashboardPageComponent },
  { path: '**', redirectTo: '' }
];

bootstrapApplication(AppComponent, {
  providers: [provideHttpClient(), provideRouter(routes)]
}).catch(err => console.error(err));
