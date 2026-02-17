import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Book } from '../../interfaces/book';
import { BooksService } from 'src/app/services/books.service';
import { ApiService } from '../../api.service';

@Component({
  selector: 'app-book',
  standalone: true,
  imports: [CommonModule, RouterLink, FormsModule],
  templateUrl: './book.component.html',
  styleUrl: './book.component.scss'
})
export class BookComponent implements OnInit {
  book: Book | null = null;
  loading = false;
  error = '';
  bookId = '';

  // Reviews
  reviews: any[] = [];
  averageRating: number | null = null;
  totalReviews = 0;
  myReview: any = null;
  reviewForm = { rating: 5, title: '', body: '' };
  showReviewForm = false;
  reviewSubmitting = false;
  reviewMessage = '';

  constructor(
    private booksService: BooksService,
    private api: ApiService,
    private route: ActivatedRoute
  ) {}

  ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.bookId = id;
      this.fetchBook(id);
      this.loadReviews(id);
      this.loadMyReview(id);
    }
  }

  fetchBook(id: string) {
    this.loading = true;
    this.error = '';
    this.api.getPublicBook(id).subscribe({
      next: res => {
        const payload = res.data?.book || {};
        const fallback = this.booksService.getDummyBooks(1)[0];
        this.book = {
          title: payload.title || fallback.title,
          slag: payload.slug || payload.id || fallback.slag,
          cover: payload.coverImageUrl || fallback.cover,
          short_description: payload.subtitle || payload.description,
          long_description: payload.description || fallback.long_description,
          author_avatar: undefined,
          author_name: payload.authorName || 'Unknown Author',
          author_short_bio: '',
          description: payload.description,
          price: payload.priceAmount ?? undefined,
          priceCurrency: payload.priceCurrency,
          editors: [],
          paperback: ''
        };
        this.loading = false;
      },
      error: () => {
        this.book = this.booksService.getDummyBooks(1)[0];
        this.error = 'Unable to load book; displaying sample preview.';
        this.loading = false;
      }
    });
  }

  loadReviews(id: string) {
    this.api.getBookReviews(id).subscribe({
      next: res => {
        this.reviews = res.data?.reviews || [];
        this.averageRating = res.data?.averageRating ?? null;
        this.totalReviews = res.data?.totalCount || 0;
      },
      error: () => {}
    });
  }

  loadMyReview(id: string) {
    this.api.getMyReview(id).subscribe({
      next: res => {
        if (res.data) {
          this.myReview = res.data;
          this.reviewForm = {
            rating: res.data.rating,
            title: res.data.title || '',
            body: res.data.body || ''
          };
        }
      },
      error: () => {}
    });
  }

  submitReview() {
    if (!this.bookId) return;
    this.reviewSubmitting = true;
    this.reviewMessage = '';
    this.api.submitReview(this.bookId, {
      rating: this.reviewForm.rating,
      title: this.reviewForm.title || undefined,
      body: this.reviewForm.body || undefined
    }).subscribe({
      next: () => {
        this.reviewSubmitting = false;
        this.reviewMessage = 'Review saved!';
        this.showReviewForm = false;
        this.loadReviews(this.bookId);
        this.loadMyReview(this.bookId);
      },
      error: () => {
        this.reviewSubmitting = false;
        this.reviewMessage = 'Failed to save review.';
      }
    });
  }

  deleteReview() {
    if (!this.bookId) return;
    this.api.deleteReview(this.bookId).subscribe({
      next: () => {
        this.myReview = null;
        this.reviewForm = { rating: 5, title: '', body: '' };
        this.loadReviews(this.bookId);
      }
    });
  }

  starsArray(n: number): number[] {
    return Array.from({ length: n }, (_, i) => i + 1);
  }

  setRating(n: number) {
    this.reviewForm.rating = n;
  }
}
