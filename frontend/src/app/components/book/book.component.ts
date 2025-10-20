import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { Book } from '../../interfaces/book';
import { BooksService } from 'src/app/services/books.service';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-book',
  standalone: true,
  imports: [CommonModule, RouterLink],
  templateUrl: './book.component.html',
  styleUrl: './book.component.scss'
})
export class BookComponent {
  book: Book | null = null;

  constructor(
    private booksService: BooksService,
    private route: ActivatedRoute
  ) { }

  ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.book = this.booksService.getDummyBookBySlag(id);
    }
  }
}
