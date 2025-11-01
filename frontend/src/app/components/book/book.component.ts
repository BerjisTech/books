import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { CommonModule } from '@angular/common';
import { Book } from '../../interfaces/book';
import { BooksService } from 'src/app/services/books.service';
import { ApiService } from '../../api.service';

@Component({
  selector: 'app-book',
  standalone: true,
  imports: [CommonModule, RouterLink],
  templateUrl: './book.component.html',
  styleUrl: './book.component.scss'
})
export class BookComponent implements OnInit {
  book: Book | null = null;
  loading = false;
  error = '';

  constructor(
    private booksService: BooksService,
    private api: ApiService,
    private route: ActivatedRoute
  ) {}

  ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.fetchBook(id);
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
}
